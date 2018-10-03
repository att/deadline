package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnTimeEvent(test *testing.T) {
	nowEvent := getOldEvent(0)
	constraints, _ := FromBlueprint(time.Now().Add(-2*time.Hour), makeConstraint("3h"))

	assert.True(test, nowEvent.IsSuccessful(constraints))
}

func TestLateEvent(test *testing.T) {
	nowEvent := getOldEvent(0)
	constraints, _ := FromBlueprint(time.Now().Add(-2*time.Hour), makeConstraint("1h"))

	assert.False(test, nowEvent.IsSuccessful(constraints))
}

func getOldEvent(secondsAgo int64) Event {
	past := time.Now().Add(time.Duration(secondsAgo)).Unix()

	return Event{
		Name:       "test",
		ReceivedAt: past,
		Details:    make(map[string]string),
	}
}

func makeConstraint(recieveBy string) EventConstraintsBlueprint {
	return EventConstraintsBlueprint{
		ReceiveBy: recieveBy,
	}
}
