package schedule

import (
	"log"
	"testing"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/notifier"
	"github.com/stretchr/testify/assert"
)

var m = NewManager()
var c = config.Config{
	DAO: "file",
	FileConfig: config.FileConfig{
		Directory: "../server/",
	},
	DBConfig: config.DBConfig{
		ConnectionString: "somethintoo",
	},
	Server: config.ServConfig{
		Port: "8081",
	},

} 
var e1 = common.Event{
	Name:      "first event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "19:00:00",
	Success:   true,
}
var e2 = common.Event{
	Name:      "second event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "17:00:00",
	Success:   true,
}

var e3 = common.Event{
	Name:      "third event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "18:00:00",
	Success:   true,
}

var s1 = Schedule{
	Name: "First Schedule",
	Start: Node{
		Nodes: []Node{
			Node{

				Event: &e1,
			},
			Node{
				Event: &e3,
			},
		},
	},
	Handler: Handler{
			Name: "WEBHOOK",
			Address: "http://localhost:8081/api/v1/msg",
	},

}

var s2 = Schedule{
	Name: "Second Schedule",
	Start: Node{
		Nodes: []Node{
			Node{

				Event: &e1,
			},
			Node{
				Event: &e2,
			},
		},
	},
}

var s3 = Schedule{
	Name: "Third Schedule",
	Start: Node{
		Nodes: []Node{
			Node{

				Event: &e2,
			},
		},
	},
}

var fd = NewScheduleDAO(&c)
var s = Schedule{
	Name:   "sample_schedule",
	Timing: "daily",
	Handler: Handler{
		Name:    "email handler",
		Address: "kp755d@att.com",
	},
}

func TestSendFile(test *testing.T) {
	assert.Nil(test, fd.Save(s), "Could not save the file.")
}

func TestGetFile(test *testing.T) {
	f, err := fd.GetByName("sample_schedule")
	assert.Nil(test, err, "Could not find the file.")
	assert.NotNil(test, f, "Could not find the file.")

}

func TestManager(test *testing.T) {

	UpdateSchedule(m, &s1)
	UpdateSchedule(m, &s2)
	UpdateSchedule(m, &s3)
	log.Println("Current map: ", m.subscriptionTable)
	UpdateEvents(m, &e1, fd)
	UpdateEvents(m, &e2, fd)
}

var f1 = common.Event{
	Name:      "first event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "19:00:00",
	Success:   true,
}
var f2 = common.Event{
	Name:      "second event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "17:00:00",
	Success:   true,
}

var f3 = common.Event{
	Name:      "third event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "18:00:00",
	Success:   true,
}

func TestEvaluation(test *testing.T) {
	assert.False(test, EvaluateEvent(&f1,notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "It is coming back as true")
	assert.True(test, EvaluateEvent(&f2,notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as false")
	assert.False(test, EvaluateEvent(&f3,notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as true")

}
