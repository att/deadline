package common

import (
	"errors"
	"time"
)

type Event struct {
	Name       string            `json:"name" xml:"name,attr" db:"name"`
	Details    map[string]string `json:"details,omitempty" db:"details"`
	ReceivedAt int64             `json:"received-at" db:"receiveat"`
}

type EventConstraints struct {
	ReceiveBy int64
}

func (e *Event) ValidateEvent() error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

// Given a set of constraints, is the event successful.
func (e *Event) IsSuccessful(c EventConstraints) bool {
	onTime := e.ReceivedAt <= c.ReceiveBy
	return onTime
}

func FromBlueprint(startAt time.Time, blueprint EventConstraintsBlueprint) (EventConstraints, error) {
	e := EventConstraints{}

	if receiveBy, err := time.ParseDuration(blueprint.ReceiveBy); err != nil {
		return e, err
	} else {
		e.ReceiveBy = startAt.Add(receiveBy).Unix()
		return e, nil
	}
}

// func (e *Event) OnTime() bool {

// 	return true
// 	// //byTime := ConvertTime(e.ReceiveBy)
// 	// atTime := ConvertTime(e.ReceiveAt)
// 	// //Debug.Println(byTime)
// 	// Debug.Println(atTime)

// 	// if atTime.IsZero() {
// 	// 	if time.Now().After(byTime) {
// 	// 		return false
// 	// 	}
// 	// 	return true

// 	// }
// 	// if atTime.Before(byTime) {
// 	// 	return true
// 	// }

// 	// return false
// }
