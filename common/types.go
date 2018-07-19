package common

import (
	"encoding/xml"
	//"github.com/Masterminds/structable"
)

type Event struct {
	XMLName   xml.Name          `xml:"event"`
	Name      string            `json:"name" xml:"name,attr" stbl:"id"`
	Success   bool              `json:"success" xml:"success" stbl:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty" stbl:"details"`
	ReceiveBy string            `xml:"receive-by,attr" stbl:"receive_by"`
}

type Schedule struct {
	XMLName xml.Name `xml:"schedule"`
	Handler Handler  `xml:"handler" json:"handler" stbl:"handler"`
	//XMLName xml.Name 	`xml:"schedule"`
	Timing   string  `xml:"timing,attr" stbl:"timing"`
	Name     string  `xml:"name,attr"  stbl:"name"`
	Schedule []Event `xml:"event" stbl:"event"`
	ReceivedEvents []Event `xml:"receivedevents" stbl:"receivedevents"` 
}

type Handler struct {
	XMLName xml.Name `xml:"handler"`
	Name    string   `xml:"name,attr" stbl:"name"`
	Address string   `xml:"address" stbl:"address"`
}

type Error struct {
	To string `xml:"to,attr" stbl:"to"`
}


