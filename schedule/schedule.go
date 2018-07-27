package schedule

import (
	"egbitbucket.dtvops.net/deadline/common"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func (s Schedule) EventOccurred(e *common.Event) {
	//loop through schedule, find event, mark it as true
	ev := s.Start.findEvent(e.Name)
	if  ev != nil {
		if makeLive(ev) != nil {
		log.Println("We were able to locate and mark the event as true.")
		s.Start.OkTo = &s.End
		}
	}
	s.Start.ErrorTo = &s.Error

}

func makeLive(e *common.Event) error{

	log.Println("Found " + e.Name)
	e.IsLive = true
	e.ReceiveAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (start Node) findEvent(name string) *common.Event {

	if start.Event != nil {
		if start.Event.Name == name {
			log.Println("Found " + start.Event.Name)
			return start.Event
		}

	} else {
		log.Println("This is a start to a schedule.")
	}

	for j := 0; j < len(start.Nodes); j++ {
		f := start.Nodes[j].findEvent(name)
		if f != nil {
			return f
		}
	}

	return nil
}

func NewManager() *scheduleManager {
	return &scheduleManager{
		subscriptionTable: make(map[string][]*Schedule),
	}
}

func (fd fileDAO) GetByName(name string) ([]byte, error) {

	file, err := os.Open(name + ".xml")
	if err != nil {

		return nil, err

	}
	defer file.Close()

	//read in the xml file
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (fd fileDAO) Save(s Schedule) error {
	str := s.Name + ".xml"
	f, err := os.Create(str)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := xml.NewEncoder(f)
	err = encoder.Encode(s)

	if err != nil {
		return err
	}
	return nil
}

func NewScheduleDAO() ScheduleDAO {
	// return some new object
	//if we have the indicator we are using a datbase, return database struct
	return &fileDAO{}
}

func UpdateEvents(m *scheduleManager, e *common.Event, fd ScheduleDAO) {
	//once you receive an event, tell every schedule that you have it by adding it to their array
	scheds := m.subscriptionTable[e.Name]
	if scheds == nil {

		log.Println("No subscribers.")
	}

	for _, sched := range scheds {
		sched.EventOccurred(e)

	}
	e1schd := m.subscriptionTable[e.Name]
	log.Println("Looking at " + e.Name)
	for i := 0; i < len(e1schd); {
		for a := 0; a < len(e1schd[i].Start.Nodes); {
			log.Println("Is " + e1schd[i].Start.Nodes[a].Event.Name + " alive?")
			log.Println(e1schd[i].Start.Nodes[a].Event.IsLive)
			a++
		}
		i++
	}
	log.Println("Onto next event.")
}

func UpdateSchedule(m *scheduleManager, s *Schedule) {
	//loop through array and subscribe to every event, and then add itself to the map for every event
	log.Println("Address of " + s.Name)
	log.Printf("%p\n", s)
	for i := 0; i < len(s.Start.Nodes); i++ {
		//subscribe to every event
		//put into map
		scheds := m.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
		if scheds == nil {
			m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []*Schedule{s}
			continue
		}
		scheds = append(scheds, s)
		m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds
	}
	//below is purely for testing
	e1schd := m.subscriptionTable["first event"]
	e2schd := m.subscriptionTable["second event"]
	log.Println("First event:")
	for i := 0; i < len(e1schd); {
		for a := 0; a < len(e1schd[i].Start.Nodes); {
			log.Println("Is " + e1schd[i].Start.Nodes[a].Event.Name + " alive?")
			log.Println(e1schd[i].Start.Nodes[a].Event.IsLive)
			a++
		}
		i++
	}
	log.Println("Second event:")
	for j := 0; j < len(e2schd); {
		for b := 0; b < len(e2schd[j].Start.Nodes); {
			log.Println("Is " + e2schd[j].Start.Nodes[b].Event.Name + " alive?")
			log.Println(e2schd[j].Start.Nodes[b].Event.IsLive)
			b++
		}
		j++
	}
}
