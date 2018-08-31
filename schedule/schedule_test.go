package schedule

// import (
// 	"log"
// 	"math/rand"
// 	"os"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/att/deadline/common"
// 	"github.com/att/deadline/config"
// 	"github.com/att/deadline/dao"
// 	"github.com/att/deadline/notifier"
// 	"github.com/stretchr/testify/assert"
// )

// var c = config.Config{
// 	DAO: "file",
// 	FileConfig: config.FileConfig{
// 		Directory: os.TempDir() + "/deadline_test/" + strconv.Itoa(rand.Int()),
// 	},
// 	DBConfig: config.DBConfig{
// 		ConnectionString: "somethintoo",
// 	},
// 	Server: config.ServConfig{
// 		Port: "8081",
// 	},
// }

// var m *ScheduleManager

// var e1 = common.Event{
// 	Name:      "first event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "19:00:00",
// 	Success:   true,
// }

// var e2 = common.Event{
// 	Name:      "second event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "17:00:00",
// 	Success:   true,
// }

// var e3 = common.Event{
// 	Name:      "third event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "18:00:00",
// 	Success:   true,
// }

// var s1 = common.Definition{
// 	Name:   "First Schedule",
// 	Timing: "24h",
// 	Start: common.Node{
// 		Nodes: []common.Node{
// 			common.Node{

// 				Event: &e1,
// 			},
// 			common.Node{
// 				Event: &e3,
// 			},
// 		},
// 	},
// 	Handler: common.Handler{
// 		Name:    "WEBHOOK",
// 		Address: "http://localhost:8081/api/v1/msg",
// 	},
// }

// var s2 = common.Definition{
// 	Name: "Second Schedule",
// 	Start: common.Node{
// 		Nodes: []common.Node{
// 			common.Node{

// 				Event: &e1,
// 			},
// 			common.Node{
// 				Event: &e2,
// 			},
// 		},
// 	},
// }

// var s3 = common.Definition{
// 	Name: "Third Schedule",
// 	Start: common.Node{
// 		Nodes: []common.Node{
// 			common.Node{

// 				Event: &e2,
// 			},
// 		},
// 	},
// }

// var fd = dao.NewScheduleDAO(&c)
// var s = common.Definition{
// 	Name:   "sample_schedule",
// 	Timing: "daily",
// 	Handler: common.Handler{
// 		Name:    "email handler",
// 		Address: "kp755d@att.com",
// 	},
// }

// func TestGetFile(test *testing.T) {
// 	f, err := fd.GetByName("sample_schedule")
// 	assert.Nil(test, err, "Could not find the file.")
// 	assert.NotNil(test, f, "Could not find the file.")

// }

// func TestManager(test *testing.T) {
// 	m = m.Init(&c)
// 	m.AddSchedule(&s1)
// 	m.AddSchedule(&s2)
// 	m.AddSchedule(&s3)
// 	log.Println("Current map: ", m.subscriptionTable)
// 	m.UpdateEvents(&e1)
// 	m.UpdateEvents(&e2)
// }

// var f1 = common.Event{
// 	Name:      "first event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "19:00:00",
// 	Success:   true,
// }
// var f2 = common.Event{
// 	Name:      "second event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "17:00:00",
// 	Success:   true,
// }
// var f3 = common.Event{
// 	Name:      "third event",
// 	ReceiveBy: "18:00:00",
// 	ReceiveAt: "18:00:00",
// 	Success:   true,
// }

// func TestEvaluation(test *testing.T) {
// 	assert.False(test, EvaluateEvent(&f1, notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "It is coming back as true")
// 	assert.True(test, EvaluateEvent(&f2, notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as false")
// 	assert.False(test, EvaluateEvent(&f3, notifier.NewNotifyHandler(s1.Handler.Name, s1.Handler.Address)), "Came back as true")

// }

// var beforereset = common.Definition{
// 	Name:   "First Schedule",
// 	Timing: "24h",
// 	Start: common.Node{
// 		Nodes: []common.Node{
// 			common.Node{

// 				Event: &common.Event{
// 					Name:      "first event",
// 					ReceiveBy: "18:00:00",
// 					ReceiveAt: "19:00:00",
// 					Success:   true,
// 				},
// 			},
// 			common.Node{
// 				Event: &common.Event{
// 					Name:      "third event",
// 					ReceiveBy: "18:00:00",
// 					ReceiveAt: "12:00:00",
// 					Success:   true,
// 				},
// 			},
// 		},
// 	},
// 	Handler: common.Handler{
// 		Name:    "WEBHOOK",
// 		Address: "http://localhost:8081/api/v1/msg",
// 	},
// }

// var afterreset = common.Definition{
// 	Name:   "First Schedule",
// 	Timing: "24h",
// 	Start: common.Node{
// 		Nodes: []common.Node{
// 			common.Node{

// 				Event: &common.Event{
// 					Name:      "first event",
// 					ReceiveBy: "18:00:00",
// 					ReceiveAt: "",
// 					Success:   true,
// 				},
// 			},
// 			common.Node{
// 				Event: &common.Event{
// 					Name:      "third event",
// 					ReceiveBy: "18:00:00",
// 					ReceiveAt: "",
// 					Success:   true,
// 				},
// 			},
// 		},
// 	},
// 	Handler: common.Handler{
// 		Name:    "WEBHOOK",
// 		Address: "http://localhost:8081/api/v1/msg",
// 	},
// }
// var n *ScheduleManager = &ScheduleManager{
// 	subscriptionTable: make(map[string][]*Live),
// 	EvaluationTime:    time.Now(),
// }

// func TestTimings(test *testing.T) {

// 	//schedule := n.subscriptionTable["first event"]
// 	//s := schedule[0]
// 	//assert.Equal(test, &afterreset,s )
// 	//bring back te assert equal, but one that functions properly with the function

// }
