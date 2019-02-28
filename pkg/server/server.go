package server

import (
	"git.poundadm.net/anachronism/xmcctl/pkg/apis/protocol/v1"
	"net"
)

type RegisteredDevice struct {
	v1.Device
	Updates chan v1.UpdateResponse
	Notifications chan v1.Notification
	Controls chan v1.ControlResponse
	Subscriptions map[v1.NotificationTag]
}

type Server struct {
	// Devices is a list of all registered devices
	Devices []v1.Device
	// DevicesByIp is a mapping of IPs, in string form, to their corresponding devices
	DevicesByIp map[int]*v1.Device
	// UDPListeners is a mapping of port numbers to corresponding connections
	UDPListeners map[int]*net.UDPConn
}

func (s *Server) RegisterDevice(device v1.Device) {}
