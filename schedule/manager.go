package schedule

import (
	"sync"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/dao"
)

var manager *ScheduleManager
var once sync.Once

// GetManagerInstance will return the singleton of the ScheduleManager object
func GetManagerInstance(cfg *config.Config) *ScheduleManager {
	once.Do(func() {
		manager = &ScheduleManager{
			db: dao.NewScheduleDAO(cfg),
		}

		manager.loadAllSchedules()
	})
	return manager
}

func (manager *ScheduleManager) loadAllSchedules() {
	manager.schedules = make(map[string]*Schedule)
	manager.subscriptionTable = make(map[string][]*Schedule)

	blueprints, err := manager.db.LoadScheduleBlueprints()
	if err != nil {
		//log that you couldn't load blueprints. return?
	}

	for _, blueprint := range blueprints {
		if schedule, err := FromBlueprint(&blueprint); err != nil {
			// log error
		} else {
			// TODO check and log duplicates entries
			manager.schedules[schedule.Name] = schedule
			for subscription := range schedule.SubscribesTo() {
				entry := manager.subscriptionTable[subscription]
				manager.subscriptionTable[subscription] = append(entry, schedule)
			}
		}
	}

	events, err := manager.db.LoadEvents()
	if err != nil {
		//log that you couldn't load events
	} else {
		for _, event := range events {
			schedules := manager.subscriptionTable[event.Name]
			for _, schedule := range schedules {
				schedule.EventOccurred(&event)
			}
		}
	}
}

// Update updates any schedule currently alive with the event that you pass in
func (manager *ScheduleManager) Update(e *com.Event) {
	scheds := manager.subscriptionTable[e.Name]

	if scheds == nil {
		com.Info.Println("No subscribers.")
	}

	for _, schedule := range scheds {
		schedule.EventOccurred(e)
	}
}

// GetBlueprint gets a blueprint for a schedule given the name of the blueprint
func (manager *ScheduleManager) GetBlueprint(name string) (*com.ScheduleBlueprint, error) {
	return manager.db.GetByName(name)
}

// AddSchedule adds the schedule to the current list of schedules. If the schedule's start time
// it will become live and the manager will start to evaluate it. Otherwise it will be scheduled
// to become live at that time
func (manager *ScheduleManager) AddSchedule(blueprint *com.ScheduleBlueprint) {

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

// GetSchedule gets the current running schedule by the given name. If it exists, it'll
// return it, if not, it will return nil.
func (manager *ScheduleManager) GetSchedule(name string) *Schedule {

	if s, ok := manager.schedules[name]; !ok {
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
