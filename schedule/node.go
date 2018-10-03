package schedule

import (
	"time"

	"github.com/sirupsen/logrus"

	com "github.com/att/deadline/common"
)

// Name returns the name of the event node.
func (node *EventNode) Name() string {
	return node.name
}

// Next returns the next nodes for an event node. Can return an empty array if
// you cannot yet move past this node.
func (node *EventNode) Next() ([]*NodeInstance, error) {
	next := make([]*NodeInstance, 0)
	var successful bool
	var pastDue = time.Now().Unix() > node.constraints.ReceiveBy
	var received = node.event != nil

	if !received && !pastDue {
		return next, nil
	}

	if received {
		successful = node.event.IsSuccessful(node.constraints)

	} else if pastDue { // not received and past due
		log.WithFields(logrus.Fields{
			"node":       node.name,
			"reason":     "event never arrived",
			"recieve-by": time.Unix(node.constraints.ReceiveBy, 0).Format(time.RFC3339),
		}).Debug("node failed")

		next = append(next, node.errorTo)
	}

	if successful {
		next = append(next, node.okTo)
	} else {
		next = append(next, node.errorTo)
	}

	return next, nil
}

// AddEvent adds an event to the EventNode
func (node *EventNode) AddEvent(e *com.Event) {
	// if node.events == nil {
	// 	node.events = make([]com.Event, 0)
	// }

	// node.events = append(node.events, e)
	node.event = e
}

// Name returns the name of the end node.
func (node *EndNode) Name() string {
	return node.name
}

// Next for an end node returns nil for both array and error
func (node *EndNode) Next() ([]*NodeInstance, error) {
	return nil, nil
}

// Name returns the name of the start node.
func (node *StartNode) Name() string {
	return "start"
}

// Next for a start node returns an array of size 1 for it's 'to' value
func (node *StartNode) Next() ([]*NodeInstance, error) {
	next := make([]*NodeInstance, 1)
	next[0] = node.to
	return next, nil
}
