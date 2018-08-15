package schedule
import (
	
//"github.com/davecgh/go-spew/spew"

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
			EvaluationTime: time.Now(),
		}
		return &manager 
		
	}
	return m

} 


func (m *ScheduleManager) Init(cfg *config.Config) *ScheduleManager{
	common.Init(os.Stdout, os.Stdout)

	n := m.getInstance()

	Fd = NewScheduleDAO(cfg) 
	schedules, err := Fd.LoadStatelessSchedules()
	if err != nil {
		common.CheckError(err)
		return n
	}

	for _, s := range schedules {
		//make sure we are not pointing to same addresses
		s.LastRun = time.Time{}
		s.MakeNodes()
		newS := s
		n.AddSchedule(&newS)

	}
/* 
	common.Debug.Println("Our subscription table in Init: ==============================================")
	spew.Dump(n.subscriptionTable)
	common.Debug.Println("======================================================================") */

	evnts,err := Fd.LoadEvents()
	common.CheckError(err)
	for _, e := range evnts {
		n.UpdateEvents(&e)
	}
	

	//load events (later)
	
	return n
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
//TODO

	for a := range m.subscriptionTable {
		s := m.subscriptionTable[a]
		for b := 0; b < len(s); b++ {
			common.Debug.Println("Looking at " + s[b].Name)
			t, err := time.ParseDuration(s[b].Timing) 
			if !s[b].LastRun.IsZero() {
				continue
			}
			
			if err != nil {
				common.Info.Println(err)
				return
			}
			dif := time.Now().Sub(m.EvaluationTime)
			if  dif >= t {
				s[b].Start.ResetEvents()
				m.EvaluationTime = time.Now()
				s[b].LastRun = time.Time{}
				continue
			}
			
			var h = notifier.NewNotifyHandler(s[b].Handler.Name,s[b].Handler.Address)
			f := s[b].Start.findEvent(a)
			if f == nil {
				common.Info.Println("Couldn't find the event in the schedule")
			} else {
				common.Debug.Println("----------------------------------------------")
				common.Debug.Println(f.Name)
				if !f.EvaluateEvent(h) {
					if f.ReceiveAt != "" {
					common.Info.Println("Just letting you know that " + f.Name + " failed!")
					s[b].LastRun = time.Now()
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
