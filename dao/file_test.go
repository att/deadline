package dao

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/att/deadline/common"
	"github.com/stretchr/testify/assert"
)

var dao = cleanAndRefreshDAO(nil, randomTempDir())

var simpleSchedule = common.Definition{
	Name:   "sample_schedule",
	Timing: "daily",
	Handler: common.Handler{
		Name:    "email handler",
		Address: "kp755d@att.com",
	},
}

var singleEventSchedule = ScheduleBlueprint{
	Timing: "daily",
	Name:   "single_event_schedule",
	Events: []EventBlueprint{
		{
			Name: "onlyEvent",
			Constraints: EventConstraints{
				ReceiveBy: "3h",
			},
		},
	},
	Handlers: []HandlerBlueprint{
		{
			Email: EmailHandlerBlueprint{
				To: "jo424n@att.com",
			},
		},
	},
	Start: StartBlueprint{
		To: "onlyEvent",
	},
	End: EndBlueprint{
		Name: "scheduleEnd",
	},
}

func TestSaveSchedule(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, randomTempDir())
	assert.Nil(test, dao.Save(&singleEventSchedule), "Could not save the file.")
}

func TestGetFile(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, "testdata/")

	blueprint, err := dao.GetByName("single_event_schedule")
	assert.Nil(test, err, "Could not find the file.")
	assert.NotNil(test, blueprint, "Could not find the file.")
	assert.Equal(test, singleEventSchedule, *blueprint, "Read file, but result is not what's expected")
}

func cleanAndRefreshDAO(dao *fileDAO, path string) *fileDAO {
	if dao == nil {
		dao = newFileDAO(path)

	} else {
		oldPath := dao.path
		if strings.HasPrefix(oldPath, os.TempDir()) {
			_ = os.RemoveAll(oldPath)
		}
	}

	return newFileDAO(path)
}

func randomTempDir() string {
	return os.TempDir() + "/deadline_test/" + strconv.Itoa(rand.Int())
}
