package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"

	com "github.com/att/deadline/common"
)

var cyclicBlueprint = &com.ScheduleBlueprint{
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

func TestSimpleSchedule(test *testing.T) {
	schedule, err := FromBlueprint(simpleBlueprint)
	assert.NotNil(test, schedule, "schedule should be not nil")
	assert.Nil(test, err, "should not have thrown an error")
}
