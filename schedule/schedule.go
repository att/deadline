package schedule

import (
	"time"
	"github.com/att/deadline/common"
	"github.com/att/deadline/notifier"
	"bytes"
	"encoding/xml"
	

)


var scheduleSchema = `
CREATE TABLE schedules (
    name text,
	timing text
)`
var scheduleEventSchema = `
CREATE TABLE schedulevents (
	schedulename text,
	ename		text,
	ereceiveby text
)`

var  eventSchema = `
CREATE TABLE events (
	name 		text,
	receiveat 	text,
	success		text,
	islive		text
)`

var handlerSchema = `
CREATE TABLE handlers (
	schedulename text,
    name text,
	address text
)


`

func ConvertTime(timing string) (time.Time){
	var m = int(time.Now().Month())
	loc, err := time.LoadLocation("Local")
	common.CheckError(err)
	parsedTime, err := time.ParseInLocation("15:04:05", timing, loc)
	if err != nil {
		parsedTime = time.Time{}
	}
	if !parsedTime.IsZero() {
		parsedTime = parsedTime.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)
	}
	return parsedTime

}


func (e *Event) EvaluateSuccess() bool {
	if (!e.IsLive) {
		return true
	}
	return e.Success
}
func (e *Event) EvaluateEvent(h notifier.NotifyHandler) bool {
	return e.EvaluateTime(h) && e.EvaluateSuccess()
}

func (s *Live) EventOccurred(e *Event) {

	ev := s.Start.findEvent(e.Name)
	
	if ev != nil {
		ev.ReceiveAt = e.ReceiveAt
		ev.IsLive = true
		ev.Success = e.Success
		s.Start.OkTo = &s.End
		
	} else {
	s.Start.ErrorTo = &s.Error
	}
	
}


func (s *Definition) MakeNodes() {
	s.fixSchedule()
	var f Event
	buf := bytes.NewBuffer(s.ScheduleContent)
				dec := xml.NewDecoder(buf)
				for dec.Decode(&f) == nil {
					e := f
					valid := e.ValidateEvent()
						if valid != nil {
							common.Debug.Println("You had an invalid event")
							return 
						}
					node1 := Node{
						Event: &e,
						Nodes: []Node{},
					}
					s.Start.Nodes = append(s.Start.Nodes, node1)
				}
}


func (s *Definition) fixSchedule() {
	evnts := []Event{}
	b := bytes.NewBuffer(s.ScheduleContent)
	d := xml.NewDecoder(b)

	for {
		t, err := d.Token()
		if err != nil  {
            break
        }

        switch et := t.(type) {

        case xml.StartElement:
            if et.Name.Local == "event" {
                c := &Event{}
                if err := d.DecodeElement(&c, &et); err != nil {
                    panic(err)
                }
				evnts = append(evnts,(*c))
            } 
		case xml.EndElement:
			break
        }
	}
	bytes, err := xml.Marshal(evnts)
	common.CheckError(err)
	s.ScheduleContent = bytes
}


func (s *Definition) ConvertToLive() *Live{


	return &Live{
		Name: s.Name,
		Timing: s.Timing,
		Handler: s.Handler,
		Start: s.Start,
	}
}