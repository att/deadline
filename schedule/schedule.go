package schedule

import (
	"encoding/xml"
	//	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"egbitbucket.dtvops.net/deadline/common"
)

func (s Schedule) EventOccurred(e common.Event) {
	//loop through schedule, find event, mark it as true
	if findEvent(s.Start, e.Name) {
		log.Println("We were able to locate and mark the event as true.")
	}

}

func findEvent(start Node, name string) bool {
	log.Print("Looking at children under " + start.Event.Name)
	if start.Event.Name == name {
		log.Println("Found " + start.Event.Name)
		start.Event.IsLive = true
		start.Event.ReceiveAt = time.Now().Format("2006-01-02 15:04:05")
		return true
	}
	for j := 0; j < len(start.Nodes); j++ {
		f := findEvent(start.Nodes[j], name)
		if f == true {
			return true
		}
	}

	log.Println("Could not find " + name + " under " + start.Event.Name)
	return false
}

func NewManager() *scheduleManager {
	return &scheduleManager{
		subscriptionTable: make(map[string][]Schedule),
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

func UpdateEvents(m *scheduleManager, e common.Event, fd ScheduleDAO) {
	//once you receive an event, tell every schedule that you have it by adding it to their array
	scheds := m.subscriptionTable[e.Name]
	if scheds == nil {

		log.Println("No subscribers.")
	}

	for _, sched := range scheds {
		sched.EventOccurred(e)
		fd.Save(sched)
	}
	log.Println("Current map: ", m.subscriptionTable)
}

func UpdateSchedule(m *scheduleManager, s Schedule) {
	//loop through array and subscribe to every event, and then add itself to the map for every event

	for i := 0; i < len(s.Start.Nodes); i++ {
		//subscribe to every event
		//put into map
		scheds := m.subscriptionTable[(s.Start.Nodes[i].Event.Name)]
		if scheds == nil {
			m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = []Schedule{s}
			continue
		}
		scheds = append(scheds, s)
		m.subscriptionTable[(s.Start.Nodes[i].Event.Name)] = scheds

	}
	log.Println("Current map: ", m.subscriptionTable)
}
