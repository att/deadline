package common

import (
	"encoding/xml"
)

type Event struct {
	XMLName   xml.Name          `xml:"event"`
	Name      string            `json:"name" xml:"name,attr"`
	Success   bool              `json:"success" xml:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty"`
	ReceiveBy string            `xml:"receive-by,attr"`
	ReceiveAt string            `xml:"receive-at,attr"`
	IsLive bool `xml:"islive"`
}

type Handler struct {
	XMLName xml.Name `xml:"handler"`
	Name    string   `xml:"name,attr"`
	Address string   `xml:"address"`
}
