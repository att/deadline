package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	com "github.com/att/deadline/common"
)

var cyclicBlueprint = &com.ScheduleBlueprint{
	Name:     "cyclicBlueprint",
	StartsAt: "2018-09-03T00:00:00+07:00",
	Timing:   "daily",
	Start: com.StartBlueprint{
		To: "firstEvent",
	},
	Events: []com.EventBlueprint{
		com.EventBlueprint{
			Name:    "firstEvent",
			OkTo:    "secondEvent",
			ErrorTo: "endNode",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "3h",
			},
		},
		com.EventBlueprint{
			Name:    "secondEvent",
			OkTo:    "firstEvent",
			ErrorTo: "endNode",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "4h",
			},
		},
	},
	End: com.EndBlueprint{
		Name: "endNode",
	},
}

var simpleBlueprint = &com.ScheduleBlueprint{
	Name:     "simpleBlueprint",
	StartsAt: "2018-09-03T00:00:00+07:00",
	Timing:   "daily",
	Start: com.StartBlueprint{
		To: "firstEvent",
	},
	Events: []com.EventBlueprint{
		com.EventBlueprint{
			Name:    "firstEvent",
			OkTo:    "secondEvent",
			ErrorTo: "endNode",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "4h",
			},
		},
		com.EventBlueprint{
			Name:    "secondEvent",
			OkTo:    "endNode",
			ErrorTo: "emailHandler",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "12h",
			},
		},
	},
	Handlers: []com.HandlerBlueprint{
		com.HandlerBlueprint{
			Name: "emailHandler",
			Email: com.EmailHandlerBlueprint{
				EmailTo: "jo424n@att.com",
			},
			To: "endNode",
		},
	},
	End: com.EndBlueprint{
		Name: "endNode",
	},
}

var hangingBlueprint = &com.ScheduleBlueprint{
	Name:     "hangingBlueprint",
	StartsAt: "2018-09-03T00:00:00+07:00",
	Timing:   "daily",
	Start: com.StartBlueprint{
		To: "firstEvent",
	},
	Events: []com.EventBlueprint{
		com.EventBlueprint{
			Name:    "firstEvent",
			OkTo:    "secondEvent",
			ErrorTo: "endNode",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "4h",
			},
		},
		com.EventBlueprint{
			Name:    "secondEvent",
			OkTo:    "endNode",
			ErrorTo: "endNode",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "4h",
			},
		},
	},
	Handlers: []com.HandlerBlueprint{
		com.HandlerBlueprint{
			Name: "emailHandler",
			Email: com.EmailHandlerBlueprint{
				EmailTo: "jo424n@att.com",
			},
			To: "endNode",
		},
	},
	End: com.EndBlueprint{
		Name: "endNode",
	},
}

func TestEmptySchedule(test *testing.T) {
	schedule, err := FromBlueprint(&com.ScheduleBlueprint{})
	assert.Nil(test, schedule, "schedule should be nil")
	assert.NotNil(test, err, "should have thrown an error")
}

func TestCyclicSchedule(test *testing.T) {
	schedule, err := FromBlueprint(cyclicBlueprint)
	assert.Nil(test, schedule, "schedule should be nil")
	assert.NotNil(test, err, "should have thrown an error")
}

func TestHangingSchedule(test *testing.T) {
	schedule, err := FromBlueprint(hangingBlueprint)
	assert.Nil(test, schedule, "schedule should be  nil")
	assert.NotNil(test, err, "should have thrown an error")
}

func TestNodeConnections(test *testing.T) {
	schedule, err := FromBlueprint(simpleBlueprint)
	assert.NotNil(test, schedule, "schedule should be not nil")
	assert.Nil(test, err, "should not have thrown an error")
	assert.NotNil(test, schedule.nodes, "node list should not be empty")
	assert.NotEmpty(test, schedule.nodes, "node list should not be empty")
	assert.Equal(test, 5, len(schedule.nodes), "nodes should be 5 items long")

	first, _ := schedule.nodes["firstEvent"]
	second, _ := schedule.nodes["secondEvent"]
	email, _ := schedule.nodes["emailHandler"]
	end := schedule.End

	assertOnEventNode(test, first, "firstEvent", second, end)
	assertOnEventNode(test, second, "secondEvent", end, email)

	assertOnHandlerNode(test, email, "emailHandler", end)
}

func TestFailedSchedule(test *testing.T) {
	schedule, err := FromBlueprint(simpleBlueprint)
	assert.NotNil(test, schedule, "schedule should be not nil")
	assert.Nil(test, err, "should not have thrown an error")

	schedule.StartTime = time.Now().Add(-23 * time.Hour)
	state := schedule.Evaluate()
	assert.Equal(test, Failed, state)

}

func TestEndedSchedule(test *testing.T) {
	schedule, err := FromBlueprint(simpleBlueprint)
	assert.NotNil(test, schedule, "schedule should be not nil")
	assert.Nil(test, err, "should not have thrown an error")

	schedule.StartTime = time.Now().Add(-24 * time.Hour)
	schedule.EventOccurred(com.Event{
		ReceivedAt: time.Now().Add(-22 * time.Hour).Unix(),
		Name:       "firstEvent",
	})
	schedule.EventOccurred(com.Event{
		ReceivedAt: time.Now().Add(-14 * time.Hour).Unix(),
		Name:       "secondEvent",
	})
	state := schedule.Evaluate()
	assert.Equal(test, "ended", state.String())

}

func assertOnEventNode(test *testing.T, node *NodeInstance, name string, okTo *NodeInstance, errTo *NodeInstance) {
	assert.NotNil(test, node, name+" node should not be nil")
	ev, ok := node.value.(EventNode)
	assert.True(test, ok, "")
	assert.Equal(test, okTo, ev.okTo, "")
	assert.Equal(test, errTo, ev.errorTo, "")
}

func assertOnHandlerNode(test *testing.T, node *NodeInstance, name string, to *NodeInstance) {
	assert.NotNil(test, node, name+" node should not be nil")
	email, ok := node.value.(EmailHandlerNode)
	assert.True(test, ok, "")
	assert.Equal(test, to, email.to, "")
}
