package schedule
import (
"time"
"log"
"egbitbucket.dtvops.net/deadline/notifier"
)


func NewManager() *ScheduleManager {

	
	return &ScheduleManager{
		subscriptionTable: make(map[string][]*Schedule),
	}
}

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

func (m *ScheduleManager) EvaluateAll(scheds []Schedule) {


		for b := 0; b < len(scheds); b++ {

			t, err := time.ParseDuration(scheds[b].Timing) 
			if err != nil {
				log.Println(err)
				return
			}
			if time.Now().Sub(m.EvaluationTime) > t  {
				var h = notifier.NewNotifyHandler(scheds[b].Handler.Name,scheds[b].Handler.Address)
				
				//go through every event
				f := scheds[b].Start.findEvent(a)
				if f == nil {
						log.Println("Couldn't find the event in the schedule")
						return
				} else {
						log.Println("----------------------------------------------")
						log.Println(f.Name)
						f.EvaluateEvent(h)	
				} 
			}
		}
	m.EvaluationTime = time.Now()

}