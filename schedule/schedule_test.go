package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"

	com "github.com/att/deadline/common"
)

var cyclicBlueprint = &com.ScheduleBlueprint{
	Name: "cyclicBlueprint",
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
	Name: "simpleBlueprint",
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

var hangingBlueprint = &com.ScheduleBlueprint{
	Name: "hangingBlueprint",
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

func TestSimpleSchedule(test *testing.T) {
	schedule, err := FromBlueprint(simpleBlueprint)
	assert.NotNil(test, schedule, "schedule should be not nil")
	assert.Nil(test, err, "should not have thrown an error")
	assert.NotNil(test, schedule.nodes, "node list should not be empty")
	assert.NotEmpty(test, schedule.nodes, "node list should not be empty")
	assert.Equal(test, 5, len(schedule.nodes), "nodes should be 5 items long")

	node, _ := schedule.nodes["firstEvent"]
	assertOnEventNode(test, node, "firstEvent")

	node, _ = schedule.nodes["secondEvent"]
	assertOnEventNode(test, node, "secondEvent")

	node, _ = schedule.nodes["emailHandler"]
	assertOnHandlerNode(test, node, "emailHandler")
}

func assertOnEventNode(test *testing.T, node *NodeInstance, name string) {
	assert.NotNil(test, node, name+" node should not be nil")
}

func assertOnHandlerNode(test *testing.T, node *NodeInstance, name string) {
	assert.NotNil(test, node, name+" node should not be nil")
}
