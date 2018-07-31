package schedule

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"egbitbucket.dtvops.net/deadline/common"
	"github.com/jasonlvhit/gocron"
)

func EvaluateTime(by string, at string) bool {
	//get string value, convert it
	//parse for the following: 15:04:05
	loc, err := time.LoadLocation("Local")
    if err != nil {
        panic(err)
    }
	
	byParse, err := time.ParseInLocation("15:04:05", by, loc)
	if err != nil {
		log.Println("Could not find receive by time")
	}
	var m = int(time.Now().Month())
	byParse = byParse.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)
	log.Println(byParse)
	atParse, err := time.ParseInLocation("15:04:05", at, loc)
	if err != nil {
		log.Println(time.Now())
		if time.Now().After(byParse) {
			//if it has passed the time, return that it's late
			log.Println("The event is late")
			return false
		}
		log.Println("Could not find receive at time, but the time has not come yet.")
		return true

	}	
	atParse = atParse.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)
	if atParse.Before(byParse){
		log.Println("The event is here and it is not late!")
		return true
	}
	log.Println("The event was received, but it was late.")
	return false
}
func EvaluateSuccess(e *common.Event) bool {
	return e.Success
}
func EvaluateEvent(e *common.Event) bool {
	//want to check receive by  and receive at
	if (EvaluateTime(e.ReceiveBy, e.ReceiveAt) == true) && (EvaluateSuccess(e) == true) {
		return true
	}
	return false
}

func (s Schedule) EventOccurred(e *common.Event) {
	//loop through schedule, find event, mark it as true
	ev := s.Start.findEvent(e.Name)
	if ev != nil {
		if makeLive(ev) != nil {
			log.Println("We were able to locate and mark the event as true.")
			s.Start.OkTo = &s.End
		}
	}
	s.Start.ErrorTo = &s.Error

}

func EvaluateAll(m *scheduleManager) {
	log.Println("We are evaluating")
	for a := range m.subscriptionTable {
		s := m.subscriptionTable[a]
		for b := 0; b < len(s); b++ {
			log.Println("----------------------")
			f := s[b].Start.findEvent(a)
			log.Println("----------------------")
			if f == nil {
				fmt.Println("Couldn't find the event in the schedule")
				return
			} else {
				log.Println("----------------------------------------------")
				log.Println(f.Name)
				EvaluateEvent(f)
				
			}
		}

	}

}

//could be a function of this interface later

func (err Node) throwError() {
	log.Println("This event did not have success")
	//and other things that kill the event
	//log fatal? etc
}
func makeLive(e *common.Event) error {

	log.Println("Found " + e.Name)
	e.IsLive = true
	e.ReceiveAt = time.Now().Format("15:04:05")
	return nil
}
func (start Node) findEvent(name string) *common.Event {

	if start.Event != nil {
		log.Println("Checking " + name)
		if start.Event.Name == name {
			return start.Event
		}
		
	} else {
		log.Println("This is a start to a schedule.")
	}
	
	for j := 0; j < len(start.Nodes); j++ {
		f := start.Nodes[j].findEvent(name)
		if f != nil {
			log.Println("Found " + name + " in traversal for event evaluation.")
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
			if e2schd[j].Start.Nodes[b].Event.IsLive {
				log.Println(e2schd[j].Start.Nodes[b].Event.ReceiveAt)
			}
			b++
		}
		j++
	} 

}

func UpdateSchedule(m *scheduleManager, s *Schedule) {
	//loop through array and subscribe to every event, and then add itself to the map for every event
	go gocron.Every(10).Seconds().Do(EvaluateAll, m)
	go gocron.Start()
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
	/**/

}
