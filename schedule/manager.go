package schedule
import (
//"github.com/davecgh/go-spew/spew"
"errors"
"time"
"os"
"egbitbucket.dtvops.net/deadline/common"
"egbitbucket.dtvops.net/deadline/notifier"
"egbitbucket.dtvops.net/deadline/config"
)

var Fd ScheduleDAO


func (m *ScheduleManager) getInstance() *ScheduleManager{
	if m == nil {
		var manager = ScheduleManager{
			subscriptionTable: make(map[string][]*Schedule),
			scheduleTable: make(map[string]*Schedule),
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
	schedules, err := Fd.LoadStatelessSchedules()
	if err != nil {
		common.CheckError(err)
		return currentManager
	}

	for _, s := range schedules {
		
		s.LastRun = time.Time{}
		s.MakeNodes()
		newSchedule := s
		//make sure pointers are different
		currentManager.AddSchedule(&newSchedule)
		currentManager.scheduleTable[s.Name] = &newSchedule

	}

	evnts,err := Fd.LoadEvents()
	common.CheckError(err)
	for _, e := range evnts {
		currentManager.UpdateEvents(&e)
	}
/* 	common.Debug.Println("Our scheduleTable:")
	spew.Dump(currentManager.scheduleTable) */
	return currentManager
}

func (m *ScheduleManager) UpdateEvents(e *Event) {
	scheds := m.subscriptionTable[e.Name]
	
	if scheds == nil {

		common.Info.Println("No subscribers.")
	}
	for _, sched := range scheds {
		sched.EventOccurred(e)
	}
	
}

func (m *ScheduleManager) UpdateSchedule(s *Schedule) {
	err := Fd.Save(s)
	common.CheckError(err)
	m.AddSchedule(s)
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
					if f.ReceiveAt != "" {
					common.Info.Println(f.Name + " failed!")
					schedules[s].LastRun = time.Now()
					}
				}
				
			}
		}

	}

}

func (m *ScheduleManager) AddSchedule(s *Schedule) {

	
	for i := 0; i < len(s.Start.Nodes); i++ {
		scheds := m.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
		if scheds == nil {
			m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []*Schedule{s}
			continue
		}
		scheds = append(scheds, s)
		m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds
		
	}
	
}

func (m *ScheduleManager) GetLiveSchedule(name string) ([]byte,error) {
	
	//check if it is in the table 
	i, ok := m.scheduleTable[name]
	if !ok {
		return []byte{},errors.New("It is not in the table")
	}

	bytes, err := retrieveLiveSchedule(*i) 
	if err != nil {
		return []byte{},errors.New("Could not encode a live schedule")
	}
	//if it is, put it into the live struct 

	return bytes,nil
}