package remote

import (
	"context"
	"encoding/xml"
	"git.poundadm.net/anachronism/xmcctl/pkg/apis/v1"
	log "github.com/sirupsen/logrus"
	"net"
)

type EmotivaRemote struct {
	Name          string
	Model         string
	ControlAddr   net.UDPAddr
	NotifyAddr    net.UDPAddr
	InfoAddr      net.UDPAddr
	Subscriptions []string
}

func NewEmotivaRemote(name string, model string, controlAddr net.UDPAddr, notifyAddr net.UDPAddr, infoAddr net.UDPAddr) *EmotivaRemote {
	return &EmotivaRemote{
		Name:        name,
		Model:       model,
		ControlAddr: controlAddr,
		NotifyAddr:  notifyAddr,
		InfoAddr:    infoAddr,
	}
}

func NewEmotivaRemoteFromTransponderResponse(addr *net.UDPAddr, tr v1.TransponderResponse) *EmotivaRemote {
	r := NewEmotivaRemote(
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

// DiscoverTransponders broadcasts a discovery packet and listens for responses on the
// passed bindAddr
func DiscoverTransponders(ctx context.Context, bindAddr *net.UDPAddr) []*EmotivaRemote {
	// bind to the port identity responses will be sent to
	listener, err := net.ListenUDP("udp", bindAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"addr": bindAddr,
			"err":  err,
		}).Error("unable to bind to discovery response port")
	}

	// setup a routine to loop over all responses received by the listener
	remotes := make(chan *EmotivaRemote, 10)
	go func() {
		respPacket := make([]byte, 0, 2048)
		for {
			// if our parent context was closed, exit the read loop
			select {
			case <-ctx.Done():
				break
			default:
			}

			// we'll need the remote address from the packets to build the remotes
			_, raddr, err := listener.ReadFromUDP(respPacket)
			if err != nil {
				log.WithFields(log.Fields{
					"addr": raddr,
					"err":  err,
				}).Error("error reading ident response packet")
			}

			var tr v1.TransponderResponse
			err = xml.Unmarshal(respPacket, tr)
			if err != nil {
				// if we encounter a decoding error, log it, reset the buffer, and loop again
				log.WithFields(log.Fields{
					"packet": string(respPacket),
					"err":    err,
				}).Error("error decoding response packet")
				respPacket = respPacket[0:0]
				continue
			}
			log.WithFields(log.Fields{
				"addr":       raddr.String(),
				"model":      tr.Model,
				"name":       tr.Name,
				"apiVersion": tr.Control.Version,
			}).Info("found transponder")
			// otherwise, just clear the packet buffer
			respPacket = respPacket[0:0]

			remote := NewEmotivaRemoteFromTransponderResponse(
				raddr,
				tr,
			)
			remotes <- remote
		}
	}()

	found := make([]*EmotivaRemote, 0)
	select {
	case <-ctx.Done():
		log.Info("discovery timeout reached")
	case r := <-remotes:
		found = append(found, r)
	}
	return found
}
