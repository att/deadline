package schedule

import (
	//	"log"
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

	UpdateSchedule(m, &s1)
	UpdateSchedule(m, &s2)
	UpdateSchedule(m, &s3)
	log.Println("Current map: ", m.subscriptionTable)
	UpdateEvents(m, &e1, fd)
	UpdateEvents(m, &e2, fd)

	//these are not test cases, just here for the purpose of seeing output
	log.Println("---------------------------------------------")
	log.Println("Here lies the subscription map:")
	e1schd := m.subscriptionTable["first event"]
	e2schd := m.subscriptionTable["second event"]
	e3schd := m.subscriptionTable["third event"]
	//gigantic inefficient loops that are not good, pls fix later
	log.Println("First event:")
	for i := 0; i < len(e1schd); {
		for a := 0; a < len(e1schd[i].Start.Nodes); {
			log.Println("Is " + e1schd[i].Start.Nodes[a].Event.Name + " alive?")
			log.Println(e1schd[i].Start.Nodes[a].Event.IsLive)
			a++
		}
		i++
	}
	log.Println("Second event:")
	for j := 0; j < len(e2schd); {
		for b := 0; b < len(e2schd[j].Start.Nodes); {
			log.Println("Is " + e2schd[j].Start.Nodes[b].Event.Name + " alive?")
			log.Println(e2schd[j].Start.Nodes[b].Event.IsLive)
			b++
		}
		j++
	}
	log.Println("Third event:")
	for k := 0; k < len(e3schd); {
		for c := 0; c < len(e3schd[k].Start.Nodes); {
			log.Println("Is " + e3schd[k].Start.Nodes[c].Event.Name + " alive?")
			log.Println(e3schd[k].Start.Nodes[c].Event.IsLive)
			c++
		}
		k++
	}

	log.Println("---------------------------------------------")
}
