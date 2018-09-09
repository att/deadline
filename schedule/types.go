package schedule

import (
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/dao"
)

type NodeType int

const (
	EventNodeType NodeType = iota
	EndNodeType
	StartNodeType
	HandlerNodeType
)

type Schedule struct {
	Name          string `json:"name,attr,omitempty" db:"name"`
	StartTime     time.Time
	Start         *NodeInstance `json:"-"`
	End           *NodeInstance `json:"-"`
	nodes         map[string]*NodeInstance
	blueprintMaps com.BlueprintMaps
	subscribesTo  map[string]bool
}

type ScheduledHandler struct {
	ScheduleName string `db:"schedulename"`
	Name         string `db:"name"`
	Address      string `db:"address"`
}

type ScheduleManager struct {
	subscriptionTable map[string][]*Schedule
	schedules         map[string]*Schedule
	db                dao.ScheduleDAO
}

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
