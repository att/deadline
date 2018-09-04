package schedule

import (
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/dao"
)

var Fd dao.ScheduleDAO
var m *ScheduleManager

func GetManagerInstance() *ScheduleManager {
	if m == nil {
		m = &ScheduleManager{
			subscriptionTable: make(map[string][]*Schedule),
			ScheduleTable:     make(map[string]*Schedule),
			EvaluationTime:    time.Now(),
		}
	}
	return m

}

// func (m *ScheduleManager) Init(cfg *config.Config) *ScheduleManager {

// 	currentManager := GetManagerInstance()

// 	Fd = dao.NewScheduleDAO(cfg)
// 	blueprints, err := Fd.LoadScheduleBlueprints()
// 	if err != nil {
// 		common.CheckError(err)
// 		return currentManager
// 	}

// 	for _, bprint := range blueprints {

// 		newSchedule := FromBlueprint(&bprint)

// 		//newSchedule.LastRun = time.Time{}

// 		//make sure pointers are different
// 		currentManager.updateSubscriptions(newSchedule)

// 	}

// 	evnts, err := Fd.LoadEvents()
// 	common.CheckError(err)
// 	for _, e := range evnts {
// 		currentManager.UpdateEvents(&e)
// 	}
// 	return currentManager
// }

func (m *ScheduleManager) Update(e *com.Event) {
	scheds := m.subscriptionTable[e.Name]

	if scheds == nil {
		com.Info.Println("No subscribers.")
	}
	// for s := 0; s < len(scheds); s++ {
	// 	scheds[s].EventOccurred(e)

	// }

}

func (m *ScheduleManager) AddSchedule(blueprint *com.ScheduleBlueprint) {

	// 	err := Fd.Save(blueprint)
	// 	common.CheckError(err)
	// 	sched := FromBlueprint(blueprint)
	// 	m.updateSubscriptions(sched)
	// }

	// func (m *ScheduleManager) EvaluateAll() {

	// 	// for subs := range m.subscriptionTable {
	// 	// 	schedules := m.subscriptionTable[subs]
	// 	// 	for s := 0; s < len(schedules); s++ {

	// 	// 		t, err := time.ParseDuration(schedules[s].Timing)
	// 	// 		if !schedules[s].LastRun.IsZero() {
	// 	// 			continue
	// 	// 		}

	// 	// 		if err != nil {
	// 	// 			common.Info.Println(err)
	// 	// 			return
	// 	// 		}
	// 	// 		dif := time.Now().Sub(m.EvaluationTime)
	// 	// 		if dif >= t {
	// 	// 			ResetEvents(&schedules[s].Start)
	// 	// 			m.EvaluationTime = time.Now()
	// 	// 			schedules[s].LastRun = time.Time{}
	// 	// 			continue
	// 	// 		}

	// 	// 		var h = notifier.NewNotifyHandler(schedules[s].Handler.Name, schedules[s].Handler.Address)
	// 	// 		f := findEvent(schedules[s].Start, subs)
	// 	// 		if f == nil {
	// 	// 			common.Info.Println("Couldn't find the event in the schedule")
	// 	// 		} else {
	// 	// 			common.Debug.Println("----------------------------------------------")
	// 	// 			common.Debug.Println(f.Name)
	// 	// 			if !EvaluateEvent(f, h) {
	// 	// 				common.Info.Println(f.Name + " failed!")
	// 	// 				schedules[s].LastRun = time.Now()
	// 	// 			}

	// 	// 		}
	// 	// 	}

	// 	// }

}

// func (m *ScheduleManager) updateSubscriptions(s *Schedule) {

// 	// for i := 0; i < len(s.Start.Nodes); i++ {
// 	// 	scheds := m.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
// 	// 	if scheds == nil {
// 	// 		m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []*Schedule{s}
// 	// 		continue
// 	// 	}
// 	// 	scheds = append(scheds, s)
// 	// 	m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds

// 	// }
// 	// m.ScheduleTable[s.Name] = s
// }

func (m *ScheduleManager) GetSchedule(name string) *Schedule {

	if s, ok := m.ScheduleTable[name]; !ok {
		return nil
	} else {
		return s
	}

	// s.Events = []common.Event{}
	//clears out old live data
	// for n := 0; n < len(s.Start.Nodes); n++ {
	// 	s.Events = append(s.Events, *(s.Start.Nodes[n].Event))
	// }
}
