package schedule

import (
	"time"

	com "github.com/att/deadline/common"
)

// Name returns the name of the event node.
func (node EventNode) Name() string {
	return node.name
}

// Next returns the next nodes for an event node. Can return an empty array if
// you cannot yet move past this node.
func (node EventNode) Next() ([]*NodeInstance, error) {

	next := make([]*NodeInstance, 0)

	// try to validate the events
	for _, event := range node.events {
		if event.IsSuccessful(node.constraints) {
			return append(next, node.okTo), nil
		}
	}

	if time.Now().Unix() > node.constraints.ReceiveBy {
		return append(next, node.errorTo), nil
	}

	return next, nil
}

// AddEvent adds an event to the EventNode
func (node *EventNode) AddEvent(e *com.Event) {
	if node.events == nil {
		node.events = make([]*com.Event, 0)
	}

	node.events = append(node.events, e)
}

// Name returns the name of the end node.
func (node EndNode) Name() string {
	return node.name
}

// Next for an end node returns nil for both array and error
func (node EndNode) Next() ([]*NodeInstance, error) {
	return nil, nil
}

// Name returns the name of the start node.
func (node StartNode) Name() string {
	return "start"
}

// Next for a start node returns an array of size 1 for it's 'to' value
func (node StartNode) Next() ([]*NodeInstance, error) {
	next := make([]*NodeInstance, 1)
	next[0] = node.to
	return next, nil
}
