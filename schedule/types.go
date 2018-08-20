package schedule

import (
	"encoding/xml"
	"time"
)

type Definition struct {
	XMLName  xml.Name       	`xml:"schedule"`
	Handler  Handler 			`json:"handler" xml:"handler,omitempty" db:"handler"`
	Timing   string         	`xml:"timing,attr,omitempty" db:"timing"`	
	Name     string         	`xml:"name,attr,omitempty" db:"name"`
	Schedule []byte         	`xml:",innerxml"`
	LastRun	 time.Time			`xml:"-"`
	Start    Node           	`xml:"-"`
	End      Node           	`xml:"-"`
	Error    Node           	`xml:"-"`
}

type Live struct {

	Timing   string         	`json:"timing,attr,omitempty" db:"timing"`	
	Name     string         	`json:"name,attr,omitempty" db:"name"`
	LastRun	 time.Time			`json:"lastrun"`
	Events []Event				`json:"events"`
	Handler  Handler			`json:"handler"`
}

type Event struct {
	XMLName   xml.Name          `json:"-" xml:"event"`
	Name      string            `json:"name" xml:"name,attr" db:"name"`
	Success   bool              `json:"success" xml:"success" db:"success"`
	Details   map[string]string `json:"details,omitempty" xml:"details,omitempty" db:"details"`
	ReceiveBy string            `json:"receiveby" xml:"receiveby,attr" db:"receiveby"`
	ReceiveAt string            `json:"receiveat" xml:"receiveat,attr" db:"receiveat"`
	IsLive bool 				`json:"islive" xml:"islive"`
}

type ScheduledEvent struct {
	ScheduleName 	string		`db:"schedulename"`
	EName			string  	`db:"ename"` 
	EReceiveBy		string  	`db:"ereceiveby"` 

}

type Handler struct {
	XMLName xml.Name 			`json:"-" xml:"handler"`
	Name    string   			`json:"name" xml:"name,attr" db:"name"`
	Address string   			`json:"address" xml:"address" db:"address"`
}

type ScheduledHandler struct {
	ScheduleName 	string	`db:"schedulename"`
	Name			string  `db:"name"` 
	Address			string  `db:"address"` 

}

type Node struct {
	Event   *Event `xml:"event"`
	Nodes   []Node        `xml:"-"`
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
	subscriptionTable map[string][]*Definition
	ScheduleTable 	  map[string]*Definition
	EvaluationTime	time.Time
}

type Error struct {
	To string `xml:"to,attr"`
}

type ScheduleDAO interface {
	GetByName(string) ([]byte, error)
	Save(s *Definition) error
	LoadStatelessSchedules() ([]Definition,error)
	LoadEvents() ([]Event,error)
	SaveEvent( e *Event) error
}

type fileDAO struct{
	Path string

}

type dbDAO struct {
	ConnectionString string
}

type innerbytes struct {
	XMLName xml.Name `xml:"innerbytes"`
	Hander Handler 	`xml:"handler"`
	Events []Event	`xml:"event"`
}