package schedule
import (
"encoding/json"
"errors"
"time"
"os"
"github.com/att/deadline/common"
"github.com/att/deadline/notifier"
"github.com/att/deadline/config"
)

var Fd ScheduleDAO


func (m *ScheduleManager) getInstance() *ScheduleManager{
	if m == nil {
		var manager = ScheduleManager{
			subscriptionTable: make(map[string][]*Live),
			ScheduleTable: make(map[string]*Live),
			EvaluationTime: time.Now(),
		}
		return &manager 
		
	}
	return m

} 


func (m *ScheduleManager) Init(cfg *config.Config) *ScheduleManager{
	common.Init(os.Stdout, os.Stdout)

	currentManager := m.getInstance()

	Fd = NewScheduleDAO(cfg) 
	schedules, err := Fd.LoadSchedules()
	if err != nil {
		common.CheckError(err)
		return currentManager
	}

   for _, s := range schedules {
		
		s.LastRun = time.Time{}
		newSchedule := s
		//make sure pointers are different
		currentManager.AddSchedule(&newSchedule)
		

	}

	evnts,err := Fd.LoadEvents()
	common.CheckError(err)
	for _, e := range evnts {
		currentManager.UpdateEvents(&e)
	}
	return currentManager
}

func (m *ScheduleManager) UpdateEvents(e *Event) {
	scheds := m.subscriptionTable[e.Name]
	
	if scheds == nil {
		common.Info.Println("No subscribers.")
	}
	for s := 0; s < len(scheds);s++ {
		scheds[s].EventOccurred(e)
		
	}
	
}

func (m *ScheduleManager) UpdateSchedule(s *Definition) {
	err := Fd.Save(s)
	common.CheckError(err)
	l := s.ConvertToLive()
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
			if  dif >= t {
				schedules[s].Start.ResetEvents()
				m.EvaluationTime = time.Now()
				schedules[s].LastRun = time.Time{}
				continue
			}
			
			var h = notifier.NewNotifyHandler(schedules[s].Handler.Name,schedules[s].Handler.Address)
			f := schedules[s].Start.findEvent(subs)
			if f == nil {
				common.Info.Println("Couldn't find the event in the schedule")
			} else {
				common.Debug.Println("----------------------------------------------")
				common.Debug.Println(f.Name)
				if !f.EvaluateEvent(h) {
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

func (m *ScheduleManager) GetLiveSchedule(name string) ([]byte,error) {
	

	l, ok := m.ScheduleTable[name]
	if !ok {
		return []byte{},errors.New("It is not in the table")
	}
	l.Events = []Event{} 
	//clears out old live data
	for n := 0; n < len(l.Start.Nodes); n++ {
		l.Events = append(l.Events, *(l.Start.Nodes[n].Event))
	}
	
	return json.Marshal(&l)
}
