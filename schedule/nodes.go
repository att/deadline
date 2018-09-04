package schedule

func (node EventNode) Name() string {
	return node.name
}

func (node EndNode) Name() string {
	return node.name
}

func (node StartNode) Name() string {
	return "start"
}

// func findEvent(start common.Node, name string) *common.Event {
// 	if start.Event != nil {

// 		if start.Event.Name == name {
// 			return start.Event
// 		}
// 	}
// 	for j := 0; j < len(start.Nodes); j++ {
// 		f := findEvent(start.Nodes[j], name)
// 		if f != nil {
// 			return f
// 		}
// 	}

// 	return nil
// }

// func ResetEvents(start *common.Node) {

// 	if start == nil {
// 		return
// 	}

// 	if start.Event != nil {
// 		start.Event.ReceiveAt = ""

// 	}
// 	for j := 0; j < len(start.Nodes); j++ {
// 		ResetEvents(&start.Nodes[j])
// 	}

// 	return
// }

// func throwError(err common.Node) {
// 	//things that kill the event
// }
