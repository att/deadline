package schedule
import (
	"encoding/xml"
	"encoding/json"
	"egbitbucket.dtvops.net/deadline/common"
	"errors"
	"bytes"
)

func retrieveLiveSchedule(s Schedule) ([]byte,error){
	evnt := Event{}
	eventsFromSchedule := []Event{}
	var counter = 0

	buf := bytes.NewBuffer(s.Schedule)
	dec := xml.NewDecoder(buf)
	for dec.Decode(&evnt) == nil {
		common.Debug.Println("It gets stuck in this loop")
		e := evnt
		eventsFromSchedule = append(eventsFromSchedule,e)
		counter++
	}

	if counter == 0 {
		return []byte{},errors.New("Had problems unmarshalling")
	}

	l := LiveSchedule{
		Name: s.Name,
		Timing: s.Timing,
		LastRun: s.LastRun,
		Events: eventsFromSchedule,
		Handler: s.Handler,
		}


	bytes, err := json.Marshal(&l)
	if err != nil {
		return []byte{},err
	}
	return bytes, err

	
}
