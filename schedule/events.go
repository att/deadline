package schedule
import (
	"time"
	"errors"
	"egbitbucket.dtvops.net/deadline/notifier"
	"egbitbucket.dtvops.net/deadline/common"
)

func (e Event) ValidateEvent() error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}


func (e *Event) makeLive() {
	e.IsLive = true
	e.ReceiveAt = time.Now().Format("15:04:05")
}

func (e *Event) EvaluateTime(h notifier.NotifyHandler) bool {

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
	if atTime.Before(byTime){
		h.Send("The event is here and it is not late!")
		return true
	}
	return false
}