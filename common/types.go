package common

import (
	"encoding/xml"
	"time"
)

type Event struct {
	XMLName   xml.Name          `xml:"event"`
	Name      string            `json:"name" xml:"name,attr"`
	Success   bool              `json:"success" xml:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty"`
	ReceiveBy string            `xml:"receive-by,attr"`
	ReceiveAt string            `xml:"receive-at,attr"`
	//receives will have to be time values in the future 
	IsLive 	  bool 	            `xml:"islive"`
}

type Schedule struct {
	XMLName xml.Name `xml:"schedule"`
	Handler Handler  `xml:"handler,omitempty" json:"handler"`
	Timing   string  `xml:"timing,attr,omitempty"`
	Name     string  `xml:"name,attr,omitempty"`
	Schedule []Event `xml:"event,omitempty"`
}

func (s Schedule) EventOccurred(e Event) error {
//loop through schedule, find event, mark it as true 

for i := 0; i < len(s.Schedule); i++ {
	if e.Name == s.Schedule[i].Name {
		s.Schedule[i].IsLive = true
		s.Schedule[i].ReceiveAt = time.Now().Format("2006-01-02 15:04:05") 
	}

}

	return nil
}


type Handler struct {
	XMLName xml.Name `xml:"handler"`
	Name    string   `xml:"name,attr"`
	Address string   `xml:"address"`
}

type Error struct {
	To string `xml:"to,attr"`
}


