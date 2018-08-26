package common

import (
	"errors"
	"time"
)

func (e *Event) EvaluateSuccess() bool {
	return e.Success
}

func (e Event) ValidateEvent() error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

func (e *Event) OnTime() bool {

	byTime := ConvertTime(e.ReceiveBy)
	atTime := ConvertTime(e.ReceiveAt)
	Debug.Println(byTime)
	Debug.Println(atTime)

	if atTime.IsZero() {
		if time.Now().After(byTime) {
			return false
		}
		return true

	}
	if atTime.Before(byTime) {
		return true
	}

	return false
}
