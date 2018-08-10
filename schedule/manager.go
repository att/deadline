package schedule
import (
	"github.com/davecgh/go-spew/spew"
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
	schedules,err := Fd.LoadStatelessSchedules()
	if err != nil {
		common.CheckError(err)
		return n
	}
	common.Debug.Println("Our array of STATELESS schedules")
	spew.Dump(schedules)
	for _, s := range schedules {
		n = n.AddSchedule(&s)

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

	spew.Dump(m.subscriptionTable)

	for a := range m.subscriptionTable {
		s := m.subscriptionTable[a]
		for b := 0; b < len(s); b++ {

			t, err := time.ParseDuration(s[b].Timing) 
			if err != nil {
				common.Info.Println(err)
				return
			}
			if time.Now().Sub(m.EvaluationTime) == t {
				s[b].ResetSchedule()
				m.EvaluationTime = time.Now()
				continue
			}
			var h = notifier.NewNotifyHandler(s[b].Handler.Name,s[b].Handler.Address)
			f := s[b].Start.findEvent(a)
			if f == nil {
				common.Info.Println("Couldn't find the event in the schedule")
				return
			} else {
				common.Debug.Println("----------------------------------------------")
				common.Debug.Println(f.Name)
				f.EvaluateEvent(h)
				
			}
		}

	}

}

func (m *ScheduleManager) AddSchedule(s *Schedule) *ScheduleManager{

	var o *ScheduleManager = m
	for i := 0; i < len(s.Start.Nodes); i++ {
		scheds := o.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
		if scheds == nil {
			o.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []*Schedule{s}
			continue
		}
		scheds = append(scheds, s)
		o.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds
	}
	return o
}
