package schedule

import (


	"time"
	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/notifier"
	

)


var schema1 = `
CREATE TABLE schedules (
    name text,
    timing text
)`
var schema2 = `
CREATE TABLE schedulevents (
	schedulename text,
	ename		text,
	ereceiveby text
)`

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
	return e.Success
}
func (e *Event) EvaluateEvent(h notifier.NotifyHandler) bool {
	return e.EvaluateTime(h) && e.EvaluateSuccess()
}

func (s Schedule) EventOccurred(e *Event) {

	ev := s.Start.findEvent(e.Name)
	if ev != nil {
		ev.makeLive() 
		s.Start.OkTo = &s.End
		
	} else {
	s.Start.ErrorTo = &s.Error
	}
}

func (s *Schedule) ResetSchedule() {
	//for each event, set receive-at to ""
	//reset nodes 


}