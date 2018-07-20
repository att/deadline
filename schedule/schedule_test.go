package database

import (
	"log"
	"testing"

	"egbitbucket.dtvops.net/deadline/common"

	"github.com/stretchr/testify/assert"
)

//"egbitbucket.dtvops.net/deadline/common"
//create a new database
//assert


//var dbdriver string
var m = NewManager()
var e1 = common.Event{
	Name: "first event",
}
var e2 = common.Event{
        Name: "second event",
}

var e3 = common.Event{
        Name: "third event",
}

var s1 = common.Schedule {
	Name: "First Schedule",
	Schedule: []common.Event{e1, e2},

}

var s2 = common.Schedule {
	Name: "Second Schedule",
        Schedule: []common.Event{e1, e3},

}

var s3 = common.Schedule {
	Name: "Third Schedule",
        Schedule: []common.Event{e2},

}


var fd = NewScheduleDAO()
var s = common.Schedule{
	Name:   "sample_schedule",
	Timing: "daily",
	Handler: common.Handler{
		Name:    "email handler",
		Address: "kp755d@att.com",
	},
}

func TestSendFile(test *testing.T) {
	assert.Nil(test, fd.Save(s), "Could not save the file.")

	//will put file in directory

}

func TestGetFile(test *testing.T) {
	f, err := fd.GetByName("sample_schedule")
	assert.Nil(test, err, "Could not find the file.")
	assert.NotNil(test, f, "Could not find the file.")
	log.Printf("Received the following information: %#v\n", f)
	//will get sample schedule from directory

}

func TestManager(test *testing.T) {
//	updateEvents(m,e1)
//	updateEvents(m,e2)
//	updateEvents(m,e3)
	updateSchedule(m,s1)
	updateSchedule(m,s2)
	updateSchedule(m,s3)
        updateEvents(m,e1)
        updateEvents(m,e2)



	//these are not test cases, just here for the purpose of seeing output
	log.Printf("%#v\n", m.Manager)

}

