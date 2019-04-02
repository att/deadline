package notifier

import (
	"sync"

	"github.com/att/deadline/config"
)

var once sync.Once
var notifier *Notifier

// GetInstance gets the current running instance of a Notifier class
func GetInstance(cfg *config.Config) *Notifier {
	once.Do(func() {
		notifier = &Notifier{}
	})

	return notifier
}

// Notify is the main API to notify some entity with a message of some kind
func (notifier *Notifier) Notify(notification Notification) {

}
