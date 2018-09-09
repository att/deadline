package schedule

import (
	"sync"
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/dao"
)

type NodeType int
type State int

const (
	// ExpectedTimeLayout is the expected time layout for schedules when unix time (int64) isn't used
	ExpectedTimeLayout string = time.RFC3339

	// EventNodeType are types of event nodes
	EventNodeType NodeType = iota
	// EndNodeType is an 'end' node
	EndNodeType
	// StartNodeType is a 'start' node
	StartNodeType
	// HandlerNodeType is a 'handler' node
	HandlerNodeType

	// Running is the running state
	Running State = iota
	// Ended is completed with success
	Ended
	// Failed is completed with failure
	Failed
)

var (
	TimingAilias = map[string]string{
		"weekly": "168h",
		"daily":  "24h",
		"hourly": "1h",
	}
)

type Schedule struct {
	Name          string `json:"name,attr,omitempty" db:"name"`
	StartTime     time.Time
	Start         *NodeInstance `json:"-"`
	End           *NodeInstance `json:"-"`
	nodes         map[string]*NodeInstance
	blueprintMaps com.BlueprintMaps
	subscribesTo  map[string]bool
	eventLock     *sync.RWMutex
	state         State
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
	rwLock            *sync.RWMutex
	blueprints        chan com.ScheduleBlueprint
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
