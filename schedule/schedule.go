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

func ConvertTimes(by string,at string) (time.Time,time.Time){
	var m = int(time.Now().Month())
	loc, err := time.LoadLocation("Local")
	common.CheckError(err)
	byParse, err := time.ParseInLocation("15:04:05", by, loc)
	common.CheckError(err)

	byParse = byParse.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)

	atParse, err := time.ParseInLocation("15:04:05", at, loc)
	if err != nil {
		atParse = time.Time{}
	}
	if !atParse.IsZero() {
	atParse = atParse.AddDate(time.Now().Year(),m-1,time.Now().Day()-1)
	}
	return byParse,atParse
}



func EvaluateTime(by string, at string,h notifier.NotifyHandler) bool {

	byTime,atTime := ConvertTimes(by,at)
	if atTime.IsZero() {
		if time.Now().After(byTime) {
		
			h.Send("The event is late. Never arrived.")
			return false
		}
		return true

	}
	
	if atTime.Before(byTime){
		h.Send("The event is here and it is not late!")
		return true
	}
	return false
}
func EvaluateSuccess(e *Event) bool {
	return e.Success
}
func EvaluateEvent(e *Event,h notifier.NotifyHandler) bool {

	return EvaluateTime(e.ReceiveBy, e.ReceiveAt,h) && EvaluateSuccess(e)

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