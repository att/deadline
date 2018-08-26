package common

import (
	"encoding/xml"
	"time"
)

type Definition struct {
	XMLName         xml.Name `xml:"schedule"`
	Handler         Handler  `json:"handler" xml:"handler,omitempty" db:"handler"`
	Timing          string   `xml:"timing,attr,omitempty" db:"timing"`
	Name            string   `xml:"name,attr,omitempty" db:"name"`
	ScheduleContent []byte   `xml:",innerxml"`
	Start           Node     `xml:"-" json:"-"`
}

type Live struct {
	Timing  string    `json:"timing,attr,omitempty" db:"timing"`
	Name    string    `json:"name,attr,omitempty" db:"name"`
	LastRun time.Time `json:"lastrun"`
	Events  []Event   `json:"events"`
	Handler Handler   `json:"handler"`
}

type Event struct {
	XMLName   xml.Name          `json:"-" xml:"event"`
	Name      string            `json:"name" xml:"name,attr" db:"name"`
	Success   bool              `json:"success" xml:"success" db:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty" db:"details"`
	ReceiveBy string            `json:"receiveby" xml:"receiveby,attr" db:"receiveby"`
	ReceiveAt string            `json:"receiveat" xml:"receiveat,attr" db:"receiveat"`
	IsLive    bool              `json:"islive" xml:"islive"`
}

type ScheduledEvent struct {
	ScheduleName string `db:"schedulename"`
	EName        string `db:"ename"`
	EReceiveBy   string `db:"ereceiveby"`
}

type Handler struct {
	XMLName xml.Name `json:"-" xml:"handler"`
	Name    string   `json:"name" xml:"name,attr" db:"name"`
	Address string   `json:"address" xml:"address" db:"address"`
}

type Node struct {
	Event   *Event `xml:"event"`
	Nodes   []Node `xml:"-"`
	ErrorTo *Node  `xml:"-"`
	OkTo    *Node  `xml:"-"`
}

// type Start struct {
// 	Node
// }

// type End struct {
// 	Node
// }
