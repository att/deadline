package schedule

import (
	"time"

	com "github.com/att/deadline/common"
)

type NodeType int

const (
	EventNodeType NodeType = iota
	EndNodeType
	HandlerNodeType
)

type Schedule struct {
	Name          string        `json:"name,attr,omitempty" db:"name"`
	Start         *NodeInstance `json:"-"`
	End           *NodeInstance `json:"-"`
	nodes         map[string]*NodeInstance
	blueprintMaps com.BlueprintMaps
}

type ScheduledHandler struct {
	ScheduleName string `db:"schedulename"`
	Name         string `db:"name"`
	Address      string `db:"address"`
}

type ScheduleManager struct {
	subscriptionTable map[string][]*Schedule
	ScheduleTable     map[string]*Schedule
	EvaluationTime    time.Time
}

// type innerbytes struct {
// 	XMLName xml.Name       `xml:"innerbytes"`
// 	Hander  common.Handler `xml:"handler"`
// 	Events  []common.Event `xml:"event"`
// }

type Node interface {
	// Type() string
	//Next() ([]*NodeInstance, error)
	// AddEdge(node *Node) error
	Name() string
}

type Handler interface {
	Handle() error
}

type NodeInstance struct {
	NodeType NodeType
	value    Node
}

type EventNode struct {
	name        string
	constraints com.EventConstraints
	events      []*com.Event
	okTo        *NodeInstance
	errorTo     *NodeInstance
}

// type HandlerNode struct {
// 	name    string
// 	Handler Handler
// }

type StartNode struct {
	to *NodeInstance
}

type EndNode struct {
	name string
}

type EmailHandlerNode struct {
	to      *NodeInstance
	name    string
	emailTo string
}
