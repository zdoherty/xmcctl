package v1

import "encoding/xml"

const (
	SelfIdentityRequestPort  = 7000
	SelfIdentityResponsePort = 7001
)

type SelfIdentityRequest struct {
	XMLName xml.Name `xml:"emotivaPing"`
}

type SelfIdentityResponse struct {
	XMLName xml.Name                   `xml:"emotivaTransponder"`
	Model   string                     `xml:"model"`
	Name    string                     `xml:"name"`
	Control SelfIdentityReponseControl `xml:"control"`
}

type SelfIdentityReponseControl struct {
	XMLName      xml.Name `xml:"control"`
	Version      string   `xml:"version"`
	ControlPort  int      `xml:"controlPort"`
	NotifyPort   int      `xml:"notifyPort"`
	InfoPort     int      `xml:"infoPort"`
	SetupPortTCP int      `xml:"setupPortTCP"`
}

type ControlRequest struct{}

type ControlResponse struct{}

type Notification struct{}

type SubscribeRequest struct{}

type SubscribeResponse struct{}

type UnsubscribeRequest struct{}

type UnsubscribeResponse struct{}

type UpdateRequest struct{}

type UpdateResponse struct{}

type UnknownResponse struct {
	XMLName  xml.Name
	InnerXML []byte `xml:",innerxml"`
}
