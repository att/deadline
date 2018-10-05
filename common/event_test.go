package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnTimeEvent(test *testing.T) {
	nowEvent := getOldEvent(0)
	constraints, _ := FromBlueprint(time.Now().Add(-2*time.Hour), makeConstraint("3h"))
	success, reason := nowEvent.IsSuccessful(constraints)

	assert.True(test, success)
	assert.Equal(test, "", reason)
}

func TestLateEvent(test *testing.T) {
	nowEvent := getOldEvent(0)
	constraints, _ := FromBlueprint(time.Now().Add(-2*time.Hour), makeConstraint("1h"))
	success, reason := nowEvent.IsSuccessful(constraints)

	assert.False(test, success)
	assert.Equal(test, LateEvent, reason)
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
