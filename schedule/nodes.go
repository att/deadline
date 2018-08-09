package schedule
import (

)
func (start Node) findEvent(name string) *Event {
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

func (err Node) throwError() {
	// log.Println("This event did not have success")
	//and other things that kill the event
	//log fatal? etc
}