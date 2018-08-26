package schedule

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/att/deadline/common"
	"github.com/att/deadline/notifier"
)

func ConvertTime(timing string) time.Time {
	var m = int(time.Now().Month())
	loc, err := time.LoadLocation("Local")
	common.CheckError(err)
	parsedTime, err := time.ParseInLocation("15:04:05", timing, loc)
	if err != nil {
		parsedTime = time.Time{}
	}
	if !parsedTime.IsZero() {
		parsedTime = parsedTime.AddDate(time.Now().Year(), m-1, time.Now().Day()-1)
	}
	return parsedTime

}

func EvaluateSuccess(e *common.Event) bool {
	if !e.IsLive {
		return true
	}
	return e.Success
}
func EvaluateEvent(e *common.Event, h notifier.NotifyHandler) bool {
	return EvaluateTime(e, h) && EvaluateSuccess(e)
}

func (s *Live) EventOccurred(e *common.Event) {

	ev := findEvent(s.Start, e.Name)

	if ev != nil {
		ev.ReceiveAt = e.ReceiveAt
		ev.IsLive = true
		ev.Success = e.Success
		s.Start.OkTo = &s.End

	} else {
		s.Start.ErrorTo = &s.Error
	}

}

func MakeNodes(s *common.Definition) {
	fixSchedule(s)
	var f common.Event
	buf := bytes.NewBuffer(s.ScheduleContent)
	dec := xml.NewDecoder(buf)
	for dec.Decode(&f) == nil {
		e := f
		valid := e.ValidateEvent()
		if valid != nil {
			common.Debug.Println("You had an invalid event")
			return
		}
		node1 := common.Node{
			Event: &e,
			Nodes: []common.Node{},
		}
		s.Start.Nodes = append(s.Start.Nodes, node1)
	}
}

func fixSchedule(s *common.Definition) {
	evnts := []common.Event{}
	b := bytes.NewBuffer(s.ScheduleContent)
	d := xml.NewDecoder(b)

	for {
		t, err := d.Token()
		if err != nil {
			break
		}

		switch et := t.(type) {

		case xml.StartElement:
			if et.Name.Local == "event" {
				c := &common.Event{}
				if err := d.DecodeElement(&c, &et); err != nil {
					panic(err)
				}
				evnts = append(evnts, (*c))
			}
		case xml.EndElement:
			break
		}
	}
	bytes, err := xml.Marshal(evnts)
	common.CheckError(err)
	s.ScheduleContent = bytes
}

func ConvertToLive(s *common.Definition) *Live {

	return &Live{
		Name:    s.Name,
		Timing:  s.Timing,
		Handler: s.Handler,
		Start:   s.Start,
	}
}
