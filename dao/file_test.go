package dao

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/stretchr/testify/assert"
)

var dao = cleanAndRefreshDAO(nil, randomTempDir())

var singleEventSchedule = com.ScheduleBlueprint{
	Timing:   "daily",
	Name:     "single_event_schedule",
	StartsAt: "2018-09-08T00:00:00Z",
	Events: []com.EventBlueprint{
		{
			Name: "onlyEvent",
			Constraints: com.EventConstraintsBlueprint{
				ReceiveBy: "3h",
			},
			OkTo:    "scheduleEnd",
			ErrorTo: "email error",
		},
	},
	Handlers: []com.HandlerBlueprint{
		{
			Email: com.EmailHandlerBlueprint{
				EmailTo: "jo424n@att.com",
			},
			To:   "scheduleEnd",
			Name: "email error",
		},
	},
	Start: com.StartBlueprint{
		To: "onlyEvent",
	},
	End: com.EndBlueprint{
		Name: "scheduleEnd",
	},
}

func TestSaveSchedule(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, randomTempDir())
	assert.Nil(test, dao.Save(&singleEventSchedule), "Could not save the file.")
}

func TestGetSchedule(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, "testdata/")

	blueprint, err := dao.GetByName("single_event_schedule")
	assert.Nil(test, err, "Could not find the file.")
	assert.NotNil(test, blueprint, "Could not find the file.")
	assert.Equal(test, singleEventSchedule, *blueprint, "Read file, but result is not what's expected")
}

func TestGetEvents(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, "testdata/")

	events, err := dao.EventsAfter(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC))
	assert.Nil(test, err)

	done := make(chan bool)
	total := 0

	go func() {
		for {
			_, more := <-events
			if more {
				total++
			} else {
				done <- true
				return
			}
		}
	}()

	<-done
	assert.Equal(test, 3, total)
}

func TestGetAfter(test *testing.T) {
	dao = cleanAndRefreshDAO(dao, "testdata/")

	events, err := dao.EventsAfter(time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC))
	assert.Nil(test, err)

	done := make(chan bool)
	total := 0

	go func() {
		for {
			_, more := <-events
			if more {
				total++
			} else {
				done <- true
				return
			}
		}
	}()

	<-done
	assert.Equal(test, 1, total)
}

func cleanAndRefreshDAO(dao ScheduleDAO, path string) ScheduleDAO {
	if dao == nil {
		dao, _ = NewScheduleDAO(&config.DefaultConfig)

	} else {
		oldPath := config.DefaultConfig.FileConfig.Directory
		if strings.HasPrefix(oldPath, os.TempDir()) {
			_ = os.RemoveAll(oldPath)
		}
	}
	// TODO modifying the default, shouldn't happen in non-test code, but really should be
	// function to GetDefault()
	config.DefaultConfig.FileConfig.Directory = path
	dao, _ = NewScheduleDAO(&config.DefaultConfig)

	return dao
}

func randomTempDir() string {
	return os.TempDir() + "/deadline_test/" + strconv.Itoa(rand.Int())
}
