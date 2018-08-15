package schedule

import (
	"time"
	"log"
	"testing"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/notifier"
	"github.com/stretchr/testify/assert"
)


var c = config.Config{
	DAO: "file",
	FileConfig: config.FileConfig{
		Directory: "../server/testdata",
	},
	DBConfig: config.DBConfig{
		ConnectionString: "somethintoo",
	},
	Server: config.ServConfig{
		Port: "8081",
	},

} 
var m *ScheduleManager

var e1 = Event{
	Name:      "first event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "19:00:00",
	Success:   true,
}
var e2 = Event{
	Name:      "second event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "17:00:00",
	Success:   true,
}

var e3 = Event{
	Name:      "third event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "18:00:00",
	Success:   true,
}

var s1 = Schedule{
	Name: "First Schedule",
	Timing: "24h",
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
	assert.Nil(test, fd.Save(&s), "Could not save the file.")
}

func TestGetFile(test *testing.T) {
	f, err := fd.GetByName("sample_schedule")
	assert.Nil(test, err, "Could not find the file.")
	assert.NotNil(test, f, "Could not find the file.")

}

func TestManager(test *testing.T) {
	m = m.Init(&c)
	m.UpdateSchedule(&s1)
	m.UpdateSchedule(&s2)
	m.UpdateSchedule(&s3)
	log.Println("Current map: ", m.subscriptionTable)
	m.UpdateEvents(&e1)
	m.UpdateEvents(&e2)
}

var f1 = Event{
	Name:      "first event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "19:00:00",
	Success:   true,
}
var f2 = Event{
	Name:      "second event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "17:00:00",
	Success:   true,
}
var f3 = Event{
	Name:      "third event",
	ReceiveBy: "18:00:00",
	ReceiveAt: "18:00:00",
	Success:   true,
}

func TestEvaluation(test *testing.T) {
	assert.False(test, f1.EvaluateEvent(notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "It is coming back as true")
	assert.True(test, f2.EvaluateEvent(notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as false")
	assert.False(test, f3.EvaluateEvent(notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as true")

}


var beforereset = Schedule{
	Name: "First Schedule",
	Timing: "24h",
	Start: Node{
		Nodes: []Node{
			Node{

				Event: &Event{
					Name:      "first event",
					ReceiveBy: "18:00:00",
					ReceiveAt: "19:00:00",
					Success:   true,
					
				},
			},
			Node{
				Event: &Event{
					Name:      "third event",
					ReceiveBy: "18:00:00",
					ReceiveAt: "17:00:00",
					Success:   true,
					
				},
			},
		},
	},
	Handler: Handler{
			Name: "WEBHOOK",
			Address: "http://localhost:8081/api/v1/msg",
	},

}


var afterreset  = Schedule{
	Name: "First Schedule",
	Timing: "24h",
	Start: Node{
		Nodes: []Node{
			Node{

				Event: &Event{
					Name:      "first event",
					ReceiveBy: "18:00:00",
					ReceiveAt: "",
					Success:   true,
					
				},
			},
			Node{
				Event: &Event{
					Name:      "third event",
					ReceiveBy: "18:00:00",
					ReceiveAt: "",
					Success:   true,
					
				},
			},
		},
	},
	Handler: Handler{
			Name: "WEBHOOK",
			Address: "http://localhost:8081/api/v1/msg",
	},

}
var n *ScheduleManager
func TestTimings(test *testing.T) {
	n = n.Init(&c)
	var currentTime time.Time
	currentTime = time.Now()
	newTime := time.Now().AddDate(0,0,-1)
	n.EvaluationTime = currentTime
	n.AddSchedule(&beforereset)
	n.EvaluateAll()
	n.EvaluationTime = newTime
	n.EvaluateAll()
	//bring back te assert equal, but one that functions properly with the function 
	
} 