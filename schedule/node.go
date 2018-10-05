package schedule

import (
	"time"

	"github.com/sirupsen/logrus"

	com "github.com/att/deadline/common"
)

const (
	// EventNeverArrived indicates that the event never arrived
	EventNeverArrived = "event did not arrive by the specified time"
)

// Name returns the name of the event node.
func (node *EventNode) Name() string {
	return node.name
}

// Next returns the next nodes for an event node. Can return an empty array if
// you cannot yet move past this node.
func (node *EventNode) Next() ([]*NodeInstance, *Context) {
	next := make([]*NodeInstance, 0)
	var successful bool
	var failureReason string
	var ctx *Context
	var pastDue = time.Now().Unix() > node.constraints.ReceiveBy
	var received = node.event != nil

	if !received && !pastDue {
		return next, nil
	}

	if received {
		successful, failureReason = node.event.IsSuccessful(node.constraints)

	} else if pastDue { // not received and past due
		log.WithFields(logrus.Fields{
			"node":       node.name,
			"reason":     "event never arrived",
			"recieve-by": time.Unix(node.constraints.ReceiveBy, 0).Format(time.RFC3339),
		}).Debug("node failed")

		ctx = newContext(node.name, EventNeverArrived)
		next = append(next, node.errorTo)
	}

	if successful {
		next = append(next, node.okTo)
	} else {
		ctx = newContext(node.name, failureReason)
		next = append(next, node.errorTo)
	}

	return next, ctx
}

// AddEvent adds an event to the EventNode
func (node *EventNode) AddEvent(e *com.Event) {
	node.event = e
}

// Name returns the name of the end node.
func (node *EndNode) Name() string {
	return node.name
}

// Next for an end node returns nil for both array and error
func (node *EndNode) Next() ([]*NodeInstance, *Context) {
	return nil, nil
}

// Name returns the name of the start node.
func (node *StartNode) Name() string {
	return "start"
}

// Next for a start node returns an array of size 1 for it's 'to' value
func (node *StartNode) Next() ([]*NodeInstance, *Context) {
	next := make([]*NodeInstance, 1)
	next[0] = node.to
	return next, nil
}

func newContext(name string, reason string) *Context {
	return &Context{
		FailedNoded:   name,
		FailureReason: reason,
		FailureTime:   time.Now(),
	}
}
