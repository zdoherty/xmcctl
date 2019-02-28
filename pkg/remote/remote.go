package remote

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	protov1 "git.poundadm.net/anachronism/xmcctl/pkg/apis/protocol/v1"
	log "github.com/sirupsen/logrus"
	"net"
)

type Remote struct {
	Name          string
	Model         string
	ControlAddr   net.UDPAddr
	NotifyAddr    net.UDPAddr
	InfoAddr      net.UDPAddr
	Subscriptions []string
}

func NewRemote(name string, model string, controlAddr net.UDPAddr, notifyAddr net.UDPAddr, infoAddr net.UDPAddr) *Remote {
	return &Remote{
		Name:        name,
		Model:       model,
		ControlAddr: controlAddr,
		NotifyAddr:  notifyAddr,
		InfoAddr:    infoAddr,
	}
}

func NewRemoteFromTransponderResponse(addr *net.UDPAddr, tr protov1.SelfIdentityResponse) *Remote {
	r := NewRemote(
		tr.Name,
		tr.Model,
		net.UDPAddr{
			IP:   addr.IP,
			Port: tr.Control.ControlPort,
		},
		net.UDPAddr{
			IP:   addr.IP,
			Port: tr.Control.NotifyPort,
		},
		net.UDPAddr{
			IP:   addr.IP,
			Port: tr.Control.NotifyPort,
		},
	)
	return r
}

func sendDiscoveryPacket(dstAddr *net.UDPAddr) error {
	packet := bytes.NewBuffer([]byte{})
	packet.Write([]byte(xml.Header))
	data, err := xml.Marshal(protov1.SelfIdentityRequest{})
	if err != nil {
		return err
	}
	packet.Write(data)

	// Create a connection without peers to send the broadcast packet on
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	defer conn.Close()
	if err != nil {
		return err
	}

	nbytes, err := conn.WriteToUDP(packet.Bytes(), dstAddr)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"addr":   dstAddr.String(),
		"nbytes": nbytes,
		"data":   fmt.Sprintf("%s", packet.Bytes()),
	}).Debug("sent discovery packet")

	return nil
}

// DiscoverTransponders broadcasts a discovery packet and listens for responses on the passed bindAddr
func DiscoverTransponders(ctx context.Context, bind net.IP, dests []net.IP) ([]*Remote, error) {
	found := make([]*Remote, 0)
	remotes := make(chan *Remote, 10)

	// bind to the port identity responses will be sent to
	bindAddr := &net.UDPAddr{
		IP:   bind,
		Port: protov1.SelfIdentityResponsePort,
	}
	listener, err := net.ListenUDP("udp", bindAddr)
	defer listener.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"addr": bind,
			"err":  err,
		}).Error("unable to bind to discovery response port")
		return found, err
	}

	// send the discovery packet
	for _, dest := range dests {
		destAddr := &net.UDPAddr{
			IP:   dest,
			Port: protov1.SelfIdentityRequestPort,
		}
		err = sendDiscoveryPacket(destAddr)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"addr": destAddr,
			}).Error("error sending discovery packet")
			return found, err
		}
	}

	// setup a routine to loop over all responses received by the listener
	go func() {
		respPacket := make([]byte, 2048)
		log.Debug("awaiting identity responses")
		for {
			// if our parent context was closed, exit the read loop
			select {
			case <-ctx.Done():
				log.Info("context closed, discovery listener terminating")
				return
			default:
			}

			// we'll need the remote address from the packets to build the remotes
			_, raddr, err := listener.ReadFromUDP(respPacket)
			if err != nil {
				// handle errors returned by the reader based on type
				if operr, ok := err.(*net.OpError); ok {
					fs := log.Fields{
						"err":  operr,
						"errT": fmt.Sprintf("%T", operr),
					}
					if operr.Temporary() {
						// if the error indicates that it's temporary, retry
						log.WithFields(fs).Warn("temporary operror while reading discovery responses")
						continue
					} else if operr.Err.Error() == "use of closed network connection" {
						// if the nested error is due to the connection closing, we're done
						log.WithFields(fs).Debug("discovery listener connection closed")
						return
					} else {
						// log unknowns
						log.WithFields(fs).Error("unknown operror while reading discovery responses")
						return
					}
				}
				// if the error isn't a known type, also log its type for debugging and retry
				log.WithFields(log.Fields{
					"addr": raddr,
					"err":  err,
					"errt": fmt.Sprintf("%T", err),
				}).Error("unknown error reading ident response packet")
				continue
			}
			log.WithFields(log.Fields{
				"addr": raddr,
				"body": string(bytes.Trim(respPacket, "\x00")),
			}).Debug("got packet on ident response port")

			tr := protov1.SelfIdentityResponse{}
			err = xml.Unmarshal(respPacket, &tr)
			if err != nil {
				// if we encounter a decoding error, log it, reset the buffer, and loop again
				log.WithFields(log.Fields{
					"packet": string(respPacket),
					"err":    err,
					"from":   raddr,
				}).Error("error decoding response packet")
				respPacket = respPacket[0:0]
				continue
			}
			log.WithFields(log.Fields{
				"addr":       raddr.String(),
				"model":      tr.Model,
				"name":       tr.Name,
				"apiVersion": tr.Control.Version,
			}).Debug("found transponder")
			respPacket = respPacket[0:0]

			remote := NewRemoteFromTransponderResponse(
				raddr,
				tr,
			)
			remotes <- remote
		}
	}()

	// once the receiver loop is setup, process the results it sends over its channel or wait
	// for our context to be closed
	for {
		select {
		case <-ctx.Done():
			// if our context is closed, our work is done and we should return results
			log.Debug("discovery context closing")
			// close the listener connection, this will cause the receiver loop to exit
			err := listener.Close()
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Debug("error closing discovery listener?")
			}
			return found, nil
		case r := <-remotes:
			// as new results are found, copy them into our results slice
			found = append(found, r)
		}
	}
}
