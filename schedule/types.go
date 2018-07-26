package schedule

import (
	"encoding/xml"

	"egbitbucket.dtvops.net/deadline/common"
)

type Schedule struct {
	XMLName  xml.Name       `xml:"schedule"`
	Handler  common.Handler `xml:"handler,omitempty"`
	Timing   string         `xml:"timing,attr,omitempty"`
	Name     string         `xml:"name,attr,omitempty"`
	Schedule []byte         `xml:",innerxml"`
	Start    Node
	End      Node
}

type Node struct {
	Event *common.Event `xml:"event"`
	Nodes []Node        `xml:",any"`
}

type scheduleManager struct {
	subscriptionTable map[string][]*Schedule
}

type Error struct {
	To string `xml:"to,attr"`
}

type ScheduleDAO interface {
	GetByName(string) ([]byte, error)
	Save(s Schedule) error
}

type fileDAO struct{}
