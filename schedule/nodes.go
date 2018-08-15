package schedule

func (start Node) findEvent(name string) *Event {
	if start.Event != nil {
		
		if start.Event.Name == name {
			return start.Event
		}
	} 
	for j := 0; j < len(start.Nodes); j++ {
		f := start.Nodes[j].findEvent(name)
		if f != nil {		
			return f
		}
	}

	return nil
}

func (start *Node) ResetEvents() {
	
	if start == nil {
		return 
	}

	if start.Event != nil {
		start.Event.ReceiveAt = ""
		
	} 
	for j := 0; j < len(start.Nodes); j++ {
		start.Nodes[j].ResetEvents()
	}

	return 
}

func (err Node) throwError() {
	//things that kill the event
}