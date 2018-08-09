package schedule
import (
	"time"
	"errors"

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