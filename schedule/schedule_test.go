package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"

	com "github.com/att/deadline/common"
)

func TestSimpleSchedule(test *testing.T) {
	schedule, err := FromBlueprint(&com.ScheduleBlueprint{})
	assert.Nil(test, schedule, "schedule should be nil")
	assert.NotNil(test, err, "should have thrown an error")
}
