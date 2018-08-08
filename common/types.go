package common

import (
	"encoding/xml"
)

type Event struct {
	XMLName   xml.Name          `xml:"event"`
	Name      string            `json:"name" xml:"name,attr" db:"name"`
	Success   bool              `json:"success" xml:"success" db:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty" db:"details"`
	ReceiveBy string            `xml:"receive-by,attr" db:"receive-by"`
	ReceiveAt string            `xml:"receive-at,attr" db:"receive-at"`
	//receives will have to be time values in the future
	IsLive bool `xml:"islive"`
}

type Handler struct {
	XMLName xml.Name `xml:"handler"`
	Name    string   `xml:"name,attr"`
	Address string   `xml:"address"`
}