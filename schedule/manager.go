package schedule
import (
"log"
"github.com/jasonlvhit/gocron"
"egbitbucket.dtvops.net/deadline/notifier"
)


func NewManager() *ScheduleManager {

	
	return &ScheduleManager{
		subscriptionTable: make(map[string][]*Schedule),
	}
}

//make a manager function -- all below
func (m *ScheduleManager) UpdateEvents(e *Event) {
	scheds := m.subscriptionTable[e.Name]
	if scheds == nil {

		log.Println("No subscribers.")
	}
	for _, sched := range scheds {
		sched.EventOccurred(e)
	}

}

func (m *ScheduleManager) UpdateSchedule(s *Schedule) {
	
	go gocron.Every(10).Seconds().Do(EvaluateAll, m)
	go gocron.Start()

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

func EvaluateAll(m *ScheduleManager) {
	for a := range m.subscriptionTable {
		s := m.subscriptionTable[a]
		for b := 0; b < len(s); b++ {
			var h = notifier.NewNotifyHandler(s[b].Handler.Name,s[b].Handler.Address)
			f := s[b].Start.findEvent(a)
			if f == nil {
				log.Println("Couldn't find the event in the schedule")
				return
			} else {
				log.Println("----------------------------------------------")
				log.Println(f.Name)
				EvaluateEvent(f,h)
				
			}
		}

	}

}