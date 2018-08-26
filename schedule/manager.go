package schedule

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/dao"
	"github.com/att/deadline/notifier"
)

var Fd dao.ScheduleDAO

func (m *ScheduleManager) getInstance() *ScheduleManager {
	if m == nil {
		var manager = ScheduleManager{
			subscriptionTable: make(map[string][]*Live),
			ScheduleTable:     make(map[string]*Live),
			EvaluationTime:    time.Now(),
		}
		return &manager

	}
	return m

}

func (m *ScheduleManager) Init(cfg *config.Config) *ScheduleManager {

	currentManager := m.getInstance()

	Fd = dao.NewScheduleDAO(cfg)
	blueprints, err := Fd.LoadSchedules()
	if err != nil {
		common.CheckError(err)
		return currentManager
	}

	for _, bprint := range blueprints {

		newSchedule := ConvertToLive(&bprint)

		newSchedule.LastRun = time.Time{}

		//make sure pointers are different
		currentManager.AddSchedule(newSchedule)

	}

	evnts, err := Fd.LoadEvents()
	common.CheckError(err)
	for _, e := range evnts {
		currentManager.UpdateEvents(&e)
	}
	return currentManager
}

func (m *ScheduleManager) UpdateEvents(e *common.Event) {
	scheds := m.subscriptionTable[e.Name]

	if scheds == nil {
		common.Info.Println("No subscribers.")
	}
	for s := 0; s < len(scheds); s++ {
		scheds[s].EventOccurred(e)

	}

}

func (m *ScheduleManager) UpdateSchedule(s *common.Definition) {
	err := Fd.Save(s)
	common.CheckError(err)
	l := ConvertToLive(s)
	m.AddSchedule(l)
}

func (m *ScheduleManager) EvaluateAll() {

	for subs := range m.subscriptionTable {
		schedules := m.subscriptionTable[subs]
		for s := 0; s < len(schedules); s++ {

			t, err := time.ParseDuration(schedules[s].Timing)
			if !schedules[s].LastRun.IsZero() {
				continue
			}

			if err != nil {
				common.Info.Println(err)
				return
			}
			dif := time.Now().Sub(m.EvaluationTime)
			if dif >= t {
				ResetEvents(&schedules[s].Start)
				m.EvaluationTime = time.Now()
				schedules[s].LastRun = time.Time{}
				continue
			}

			var h = notifier.NewNotifyHandler(schedules[s].Handler.Name, schedules[s].Handler.Address)
			f := findEvent(schedules[s].Start, subs)
			if f == nil {
				common.Info.Println("Couldn't find the event in the schedule")
			} else {
				common.Debug.Println("----------------------------------------------")
				common.Debug.Println(f.Name)
				if !EvaluateEvent(f, h) {
					common.Info.Println(f.Name + " failed!")
					schedules[s].LastRun = time.Now()
				}

			}
		}

	}

}

func (m *ScheduleManager) AddSchedule(s *Live) {

	for i := 0; i < len(s.Start.Nodes); i++ {
		scheds := m.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
		if scheds == nil {
			m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []*Live{s}
			continue
		}
		scheds = append(scheds, s)
		m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds

	}
	m.ScheduleTable[s.Name] = s
}

func (m *ScheduleManager) GetLiveSchedule(name string) ([]byte, error) {

	l, ok := m.ScheduleTable[name]
	if !ok {
		return []byte{}, errors.New("It is not in the table")
	}
	l.Events = []common.Event{}
	//clears out old live data
	for n := 0; n < len(l.Start.Nodes); n++ {
		l.Events = append(l.Events, *(l.Start.Nodes[n].Event))
	}

	return json.Marshal(&l)
}
