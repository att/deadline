package common

import (
	"errors"
	"time"
)

// Event is the struct that represents something that has happened from some other system.
type Event struct {
	Name       string            `json:"name" xml:"name,attr" db:"name"`
	Details    map[string]string `json:"details,omitempty" db:"details"`
	ReceivedAt int64             `json:"received-at" db:"receiveat"`
}

// EventConstraints represent contstaints that are placed on an event. For example 'RecieveBy' is the
// constraint that says an event must be recieved by this time.
type EventConstraints struct {
	ReceiveBy int64
}

// ValidateEvent is a helper function that validates that the event is correctly strutured. I.e.,
// the Name field must not be empty. This returns an error if the event is not valid.
func (e *Event) ValidateEvent() error {
	if e.Name == "" {
		return errors.New("event name cannot be empty")
	}

	return nil
}

// IsSuccessful returns true if the event has met it's set of contstraints.
func (e *Event) IsSuccessful(c EventConstraints) bool {
	onTime := e.ReceivedAt <= c.ReceiveBy
	return onTime
}

// FromBlueprint returns constraints for an event based on the start time and blueprints
func FromBlueprint(startAt time.Time, blueprint EventConstraintsBlueprint) (EventConstraints, error) {
	e := EventConstraints{}

	if receiveBy, err := time.ParseDuration(blueprint.ReceiveBy); err != nil {
		return e, err
	} else {
		e.ReceiveBy = startAt.Add(receiveBy).Unix()
		return e, nil
	}
}
