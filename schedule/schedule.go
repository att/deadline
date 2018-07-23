package schedule

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"egbitbucket.dtvops.net/deadline/common"
)

func (s Schedule) EventOccurred(e common.Event) {
	//loop through schedule, find event, mark it as true

	for i := 0; i < len(s.Schedule); i++ {
		if e.Name == s.Schedule[i].Name {
			s.Schedule[i].IsLive = true
			s.Schedule[i].ReceiveAt = time.Now().Format("2006-01-02 15:04:05")
		}
	}

}

func NewManager() *scheduleManager {
	return &scheduleManager{
		subscriptionTable: make(map[string][]Schedule),
	}

}

func (fd fileDAO) GetByName(name string) ([]byte, error) {

	//	var s Schedule

	file, err := os.Open(name + ".xml")
	if err != nil {

		return nil, err

	}
	fmt.Println("We could open it!")
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

func updateEvents(m *scheduleManager, e common.Event) {
	//once you receive an event, tell every schedule that you have it by adding it to their array
	var scheds []Schedule = m.subscriptionTable[e.Name]
	if scheds == nil {

		log.Println("No subscribers.")
	}

	for _, sched := range scheds {
		sched.EventOccurred(e)
	}

}

func updateSchedule(m *scheduleManager, s Schedule) {
	//loop through array and subscribe to every event, and then add itself to the map for every event

	for i := 0; i < len(s.Schedule); i++ {
		//subscribe to every event
		//put into map
		var scheds []Schedule = m.subscriptionTable[(s.Schedule[i].Name)]
		if scheds == nil {
			m.subscriptionTable[(s.Schedule[i].Name)] = []Schedule{s}
			continue
		}
		scheds = append(scheds, s)
		m.subscriptionTable[s.Schedule[i].Name] = scheds

	}

}
