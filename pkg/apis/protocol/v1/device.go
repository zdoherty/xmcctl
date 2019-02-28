package v1

import (
	"errors"
	"git.poundadm.net/anachronism/xmcctl/pkg/apis/config"
	"net"
)

type Device struct {
	Name        string
	Model       string
	IP          net.IP
	ControlAddr net.UDPAddr
	NotifyAddr  net.UDPAddr
	InfoAddr    net.UDPAddr
	SetupAddr   net.TCPAddr
}

func NewDeviceFromSelfIdentityResponse(addr net.IP, sir *SelfIdentityResponse) *Device {
	d := &Device{
		Name:  sir.Name,
		Model: sir.Model,
		IP:    addr,
		ControlAddr: net.UDPAddr{
			IP:   addr,
			Port: sir.Control.ControlPort,
		},
		NotifyAddr: net.UDPAddr{
			IP:   addr,
			Port: sir.Control.NotifyPort,
		},
		InfoAddr: net.UDPAddr{
			IP:   addr,
			Port: sir.Control.InfoPort,
		},
		SetupAddr: net.TCPAddr{
			IP:   addr,
			Port: sir.Control.SetupPortTCP,
		},
	}
	return d
}

func NewDeviceFromRawDevice(rd *config.RawDevice) (*Device, error) {
	d := &Device{
		Name:  rd.Name,
		Model: rd.Model,
	}
	ip := net.ParseIP(rd.IP)
	if ip == nil {
		return nil, errors.New("unable to parse IP: " + rd.IP)
	}
	d.IP = ip
	d.ControlAddr = net.UDPAddr{
		IP:   ip,
		Port: rd.ControlPort,
	}
	d.NotifyAddr = net.UDPAddr{
		IP:   ip,
		Port: rd.NotifyPort,
	}
	d.InfoAddr = net.UDPAddr{
		IP:   ip,
		Port: rd.InfoPort,
	}
	d.SetupAddr = net.TCPAddr{
		IP:   ip,
		Port: rd.SetupPort,
	}
	return d, nil
}
