package v1

import "encoding/xml"

const (
	IdentDiscoveryPort = 7000
	IdentResponsePort = 7001
)

type Ping struct {
	XMLName xml.Name `xml:"emotivaPing"`
}

type TransponderResponse struct {
	XMLName xml.Name `xml:"emotivaTransponder"`
	Model string `xml:"model"`
	Name string `xml:"name"`
	Control TransponderReponseControl `xml:"control"`
}

type TransponderReponseControl struct {
	XMLName xml.Name `xml:"control"`
	Version string `xml:"version"`
	ControlPort int `xml:"controlPort"`
	NotifyPort int `xml:"notifyPort"`
	InfoPort int `xml:"infoPort"`
	SetPortTCP int `xml:"setupPortTCP"`
}

type CommandRequest struct {}

type CommandResponse struct {}

type Notification struct {}

type SubscribeRequest struct {}

type SubscribeResponse struct {}

type UnsubscribeRequest struct {}

type UnsubscribeResponse struct {}

