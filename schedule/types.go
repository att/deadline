package schedule

import (
	"sync"
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/dao"
)

// NodeType is the type of node
type NodeType int

// State is the state of a schedule, like Running
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
	// TimingAilias are aliases for timings. Example: weekly = 168h. That is the string 'weekly' is 168 hours.
	TimingAilias = map[string]string{
		"weekly": "168h",
		"daily":  "24h",
		"hourly": "1h",
	}

	// StateStringLookup is a lookuptable for the State iota which
	StateStringLookup = map[State]string{
		Running: "running",
		Ended:   "ended",
		Failed:  "failed",
	}
)

func (state State) String() string {
	if ret, found := StateStringLookup[state]; found {
		return ret
	}

	return "unknown"
}

// Schedule is the type that represents a running schedule
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

// Manager is tasked with running and maintaing all the schedules. There should only be 1 per process.
// It's tasked with the creation, destruction and evaulation of all schedules.
type Manager struct {
	subscriptionTable map[string][]*Schedule
	schedules         map[string]*Schedule
	db                dao.ScheduleDAO
	rwLock            *sync.RWMutex
	blueprints        chan com.ScheduleBlueprint
	evalTicker        *time.Ticker
}

// Context is the way to pass state between different nodes in a schedule
type Context struct {
	FailedNoded   string
	FailureReason string
	FailureTime   time.Time
}

// Node is the interface for nodes in the schedules and provides ways to see what they are and how they connect
// to other Nodes.
type Node interface {
	Next() ([]*NodeInstance, *Context)
	Name() string
}

// Handler is the interface for handlers to implement so the can handle failures in a uniform way.
type Handler interface {
	Handle(*Context)
}

// NodeInstance is the actual instance of a Node interface.
type NodeInstance struct {
	NodeType NodeType
	value    Node
}

// EventNode is the Node implementing type for an event.
type EventNode struct {
	name        string
	constraints com.EventConstraints
	event       *com.Event
	okTo        *NodeInstance
	errorTo     *NodeInstance
}

// StartNode is the Node implementing type for the start of a schedule.
type StartNode struct {
	to *NodeInstance
}

// EndNode is the Node implementing type for the end of a schedule.
type EndNode struct {
	name string
}

// EmailHandlerNode is the Node & Handler implementing type for the handler that emails failures.
type EmailHandlerNode struct {
	to      *NodeInstance
	name    string
	emailTo string
}
