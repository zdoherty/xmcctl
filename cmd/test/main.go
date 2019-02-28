package main

import (
	"git.poundadm.net/anachronism/xmcctl/pkg/apis/v1"

	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	EmotivaTransponderPort = 7000
	EmotivaIdentResponsePort = 7001
	broadcastAddr = "255.255.255.255"
)

func main() {
	// setup a udp listener on the magic port emotiva uses
	addr := net.UDPAddr{
		Port: EmotivaIdentResponsePort,
		IP: net.ParseIP("0.0.0.0"),
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("error creating listener: %v\n", err)
	}


	// send the ping packet in the background
	go func() {
		// build the ping packet
		ping := v1.Ping{}
		output, err := xml.Marshal(ping)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		var packetBuf bytes.Buffer

		w := io.Writer(&packetBuf)
		w.Write([]byte(xml.Header))
		w.Write(output)

		fmt.Println("packet to send:")
		os.Stdout.Write(packetBuf.Bytes())

		bcastAddr := &net.UDPAddr{
			Port: EmotivaTransponderPort,
			IP: net.ParseIP(broadcastAddr),
		}
		localAddr := &net.UDPAddr{
			Port: 32000,
			IP: net.ParseIP("10.69.0.16"),
		}

		conn, err := net.DialUDP("udp", localAddr, bcastAddr)
		if err != nil {
			fmt.Printf("error making broadcaster: %v\n", err)
		}
		defer conn.Close()
		conn.Write(packetBuf.Bytes())
	}()

	// listen for responses
	for {
		var response v1.TransponderResponse
		packet := make([]byte, 2048)
		_, remoteAddr, err := server.ReadFromUDP(packet)
		fmt.Printf("Got a packet from %v:\n", remoteAddr)
		if err != nil {
			fmt.Printf("error reading packet: %v\n", err)
			continue
		}
		os.Stdout.Write(packet)
		if err := xml.Unmarshal(packet, &response); err != nil {
			fmt.Printf("error unmarshaling: %v\n", err)
		}
		fmt.Printf("%+v\n", response)
	}
}
