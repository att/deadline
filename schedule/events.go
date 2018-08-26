package schedule

import (
	"errors"
	"time"

	"github.com/att/deadline/common"
	"github.com/att/deadline/notifier"
)

func ValidateEvent(e common.Event) error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

func EvaluateTime(e *common.Event, h notifier.NotifyHandler) bool {

	byTime := ConvertTime(e.ReceiveBy)
	atTime := ConvertTime(e.ReceiveAt)
	common.Debug.Println(byTime)
	common.Debug.Println(atTime)

	if atTime.IsZero() {
		if time.Now().After(byTime) {

			h.Send("The event is late. Never arrived.")
			return false
		}
		return true

	}
	if atTime.Before(byTime) {
		h.Send("The event is here and it is not late!")
		return true
	}

	return false
}
