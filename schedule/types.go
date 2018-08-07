package schedule

import (
	"encoding/xml"
	"egbitbucket.dtvops.net/deadline/common"
)

type Schedule struct {
	XMLName  xml.Name       `xml:"schedule"`
	Handler  Handler 		`xml:"handler,omitempty"`
	Timing   string         `xml:"timing,attr,omitempty" db:"timing"`	
	Name     string         `xml:"name,attr,omitempty" db:"name"`
	Schedule []byte         `xml:",innerxml"`
	Start    Node           `xml:"-"`
	End      Node           `xml:"-"`
	Error    Node           `xml:"-"`
}

type Node struct {
	Event   *common.Event `xml:"event"`
	Nodes   []Node        `xml:",any"`
	ErrorTo *Node         `xml:"-"`
	OkTo    *Node         `xml:"-"`
}

type Start struct {
	Node
}

type End struct {
	Node
}

type ScheduleManager struct {
	subscriptionTable map[string][]*Schedule
}

type Error struct {
	To string `xml:"to,attr"`
}

type Handler struct {
	Name string `xml:"name,attr"`
	Address string `xml:"address"`
}
type ScheduleDAO interface {
	GetByName(string) ([]byte, error)
	Save(s Schedule) error
}

type fileDAO struct{
	Path string

}

type dbDAO struct {
	ConnectionString string
}
