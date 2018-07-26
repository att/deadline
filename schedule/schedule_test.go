package schedule

import (
	//	"log"
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

var s1 = Schedule{
	Name: "First Schedule",
	Schedule: []byte(`<event name="first event" receive-by="16:00:00" receive-at="">
			<success>false</success>
			<islive>false</islive>
			</event>
			
			<event name="second event" receive-by="18:00:00" receive-at="">
                        <success>false</success>
                        <islive>false</islive>
                        </event>`),
}

var s2 = Schedule{
	Name: "Second Schedule",
}

var s3 = Schedule{
	Name:     "Third Schedule",
	Schedule: []common.Event{e2},
}

var fd = NewScheduleDAO()
var s = Schedule{
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
	//	log.Printf("Received the following information: %#v\n", f)
	//will get sample schedule from directory

}

func TestManager(test *testing.T) {
	//	updateEvents(m,e1)
	//	updateEvents(m,e2)
	//	updateEvents(m,e3)
	UpdateSchedule(m, s1)
	UpdateSchedule(m, s2)
	UpdateSchedule(m, s3)
	UpdateEvents(m, e1, fd)
	UpdateEvents(m, e2, fd)

	//these are not test cases, just here for the purpose of seeing output
	//log.Printf("%#v\n", m.subscriptionTable)

}
