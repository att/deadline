package schedule

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/notifier"
	"github.com/jasonlvhit/gocron"
)



func EvaluateTime(by string, at string,h notifier.NotifyHandler) bool {

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
	atParse, err := time.ParseInLocation("15:04:05", at, loc)
	if err != nil {
		if time.Now().After(byParse) {
		
			h.Send("The event is late. Never arrived.")
			return false
		}
		return true

	}	
	atParse = atParse.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)
	if atParse.Before(byParse){
		h.Send("The event is here and it is not late!")
		return true
	}
	return false
}
func EvaluateSuccess(e *common.Event) bool {
	return e.Success
}
func EvaluateEvent(e *common.Event,h notifier.NotifyHandler) bool {

	if (EvaluateTime(e.ReceiveBy, e.ReceiveAt,h) == true) && (EvaluateSuccess(e) == true) {
		return true
	}
	return false
}

func (s Schedule) EventOccurred(e *common.Event) {

	ev := s.Start.findEvent(e.Name)
	if ev != nil {
		if makeLive(ev) != nil {
			s.Start.OkTo = &s.End
		}
	}
	s.Start.ErrorTo = &s.Error

}

func EvaluateAll(m *ScheduleManager) {
	for a := range m.subscriptionTable {
		s := m.subscriptionTable[a]
		for b := 0; b < len(s); b++ {
			var h = notifier.NewNotifyHandler(s[b].Handler.Name,s[b].Handler.Address)
			f := s[b].Start.findEvent(a)
			if f == nil {
				fmt.Println("Couldn't find the event in the schedule")
				return
			} else {
				log.Println("----------------------------------------------")
				log.Println(f.Name)
				EvaluateEvent(f,h)
				
			}
		}

	}

}

//could be a function of this interface later

func (err Node) throwError() {
	// log.Println("This event did not have success")
	//and other things that kill the event
	//log fatal? etc
}
func makeLive(e *common.Event) error {
	e.IsLive = true
	e.ReceiveAt = time.Now().Format("15:04:05")
	return nil
}
func (start Node) findEvent(name string) *common.Event {
	if start.Event != nil {
		
		if start.Event.Name == name {
			return start.Event
		}
		
	} else {
		
	}
	
	for j := 0; j < len(start.Nodes); j++ {
		f := start.Nodes[j].findEvent(name)
		if f != nil {		
			return f
		}
	}

	return nil
}

func NewManager() *ScheduleManager {

	
	return &ScheduleManager{
		subscriptionTable: make(map[string][]*Schedule),
	}
}

func (fd fileDAO) GetByName(name string) ([]byte, error) {

	file, err := os.Open(name + ".xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

func NewScheduleDAO(c *config.Config) ScheduleDAO {
	if (c.DAO == "file"){
	return &fileDAO{} 
}
	return &dbDAO{}
}

func UpdateEvents(m *ScheduleManager, e *common.Event, fd ScheduleDAO) {
	scheds := m.subscriptionTable[e.Name]
	if scheds == nil {

		log.Println("No subscribers.")
	}
	for _, sched := range scheds {
		sched.EventOccurred(e)
	}

}

func UpdateSchedule(m *ScheduleManager, s *Schedule) {
	
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

func (db dbDAO) GetByName(name string) ([]byte, error) {
	return []byte{},nil
}

func (db dbDAO) Save(s Schedule) error {
	
	return nil
}