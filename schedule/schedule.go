package schedule

import (
	"errors"
	"time"

	com "github.com/att/deadline/common"
)

// import (
// 	"bytes"
// 	"encoding/xml"
// 	"time"

// 	"github.com/att/deadline/dao"

// 	"github.com/att/deadline/common"
// 	"github.com/att/deadline/notifier"
// )

// func ConvertTime(timing string) time.Time {
// 	var m = int(time.Now().Month())
// 	loc, err := time.LoadLocation("Local")
// 	common.CheckError(err)
// 	parsedTime, err := time.ParseInLocation("15:04:05", timing, loc)
// 	if err != nil {
// 		parsedTime = time.Time{}
// 	}
// 	if !parsedTime.IsZero() {
// 		parsedTime = parsedTime.AddDate(time.Now().Year(), m-1, time.Now().Day()-1)
// 	}
// 	return parsedTime

// }

// func EvaluateSuccess(e *common.Event) bool {
// 	if !e.IsLive {
// 		return true
// 	}
// 	return e.Success
// }
// func EvaluateEvent(e *common.Event, h notifier.NotifyHandler) bool {
// 	return EvaluateTime(e, h) && EvaluateSuccess(e)
// }

// func (s *Schedule) EventOccurred(e *common.Event) {

// 	ev := findEvent(s.Start, e.Name)

// 	if ev != nil {
// 		ev.ReceiveAt = e.ReceiveAt
// 		ev.IsLive = true
// 		ev.Success = e.Success
// 		s.Start.OkTo = &s.End

// 	} else {
// 		s.Start.ErrorTo = &s.Error
// 	}

// }

// func MakeNodes(s *dao.ScheduleBlueprint) {
// 	fixSchedule(s)
// 	var f common.Event
// 	buf := bytes.NewBuffer(s.ScheduleContent)
// 	dec := xml.NewDecoder(buf)
// 	for dec.Decode(&f) == nil {
// 		e := f
// 		valid := e.ValidateEvent()
// 		if valid != nil {
// 			common.Debug.Println("You had an invalid event")
// 			return
// 		}
// 		node1 := common.Node{
// 			Event: &e,
// 			Nodes: []common.Node{},
// 		}
// 		s.Start.Nodes = append(s.Start.Nodes, node1)
// 	}
// }

// func fixSchedule(s *dao.ScheduleBlueprint) {
// 	evnts := []common.Event{}
// 	b := bytes.NewBuffer(s.ScheduleContent)
// 	d := xml.NewDecoder(b)

// 	for {
// 		t, err := d.Token()
// 		if err != nil {
// 			break
// 		}

// 		switch et := t.(type) {

// 		case xml.StartElement:
// 			if et.Name.Local == "event" {
// 				c := &common.Event{}
// 				if err := d.DecodeElement(&c, &et); err != nil {
// 					panic(err)
// 				}
// 				evnts = append(evnts, (*c))
// 			}
// 		case xml.EndElement:
// 			break
// 		}
// 	}
// 	bytes, err := xml.Marshal(evnts)
// 	common.CheckError(err)
// 	s.ScheduleContent = bytes
// }

func validateBluePrint(blueprint *com.ScheduleBlueprint) error {
	return nil
}

// FromBlueprint creates a Schedule struct from a blueprint.  Errors can occur for various reasons
// like invalid business rules like an event node's error-to can only go to end or a handler
// or just general malformed blueprints like nodes having cycles or hanging nodes.
func FromBlueprint(blueprint *com.ScheduleBlueprint) (*Schedule, error) {
	maps := &com.BlueprintMaps{}
	var err error

	if maps, err = com.GetBlueprintMaps(blueprint); err != nil {
		return nil, err
	}

	schedule := &Schedule{
		nodes:         make(map[string]*NodeInstance),
		blueprintMaps: *maps,
	}

	schedule.End = &NodeInstance{
		NodeType: EndNodeType,
		value: &EndNode{
			name: blueprint.End.Name,
		},
	}

	schedule.nodes[blueprint.End.Name] = schedule.End
	visited := make(map[string]bool)

	if firstEvent, found := maps.Events[blueprint.Start.To]; !found {
		return nil, errors.New("Start node needs to point to an event Node")

	} else if err := schedule.addEventBlueprint(firstEvent, visited); err != nil {
		return nil, err
	}

	if firstEvent, found := schedule.nodes[blueprint.Start.To]; !found {
		return nil, errors.New("Schedule built, but still no first node")
	} else {
		startNode := &NodeInstance{
			NodeType: StartNodeType,
			value: &StartNode{
				to: firstEvent,
			},
		}

		schedule.nodes[startNode.value.Name()] = startNode
		schedule.Start = startNode
	}

	return schedule, nil
}

// helper function to make the ok-to node sub-tree of the given blueprint
func (schedule *Schedule) makeOKToNode(blueprint com.EventBlueprint, visited map[string]bool) error {
	if _, found := schedule.nodes[blueprint.OkTo]; !found {
		okToBlueprint, isEvent := schedule.blueprintMaps.Events[blueprint.OkTo]
		okToNode, foundOkTo := schedule.nodes[blueprint.OkTo]

		if isEvent { //okTo not already made and is an event node
			if err := schedule.addEventBlueprint(okToBlueprint, visited); err != nil {
				return err
			}
		} else if foundOkTo && okToNode.NodeType == EndNodeType {
			// do nothing, just having checked for it is enough
		} else {
			return errors.New("Events can only ok-to other events or the end node")
		}
	}
	return nil
}

// helper function to make the error-to node sub-tree of the given blueprint
func (schedule *Schedule) makeErrorToNode(blueprint com.EventBlueprint, visited map[string]bool) error {
	if _, found := schedule.nodes[blueprint.ErrorTo]; !found {
		errorToBlueprint, isEvent := schedule.blueprintMaps.Events[blueprint.ErrorTo]

		if isEvent {
			if err := schedule.addEventBlueprint(errorToBlueprint, visited); err != nil {
				return err
			}
		}

		errorToHandlerBlueprint, isHandler := schedule.blueprintMaps.Handlers[blueprint.ErrorTo]
		if isHandler {
			if err := schedule.addHandlerBlueprint(errorToHandlerBlueprint, visited); err != nil {
				return err
			}
		} else {
			// at this point it wasn't found, and it wasn't an event and it wasn't a handler
			return errors.New("Couldn't find the error-to node for " + blueprint.Name)
		}
	}

	return nil
}

func (schedule *Schedule) addEventBlueprint(blueprint com.EventBlueprint, visited map[string]bool) error {
	if c, err := com.FromBlueprint(time.Now(), blueprint.Constraints); err != nil {
		return err
	} else if visited[blueprint.Name] {
		return errors.New("Possible cycle, already visited " + blueprint.Name)
	} else {

		visited[blueprint.Name] = true
		if err := schedule.makeOKToNode(blueprint, visited); err != nil {
			return err
		}

		if err := schedule.makeErrorToNode(blueprint, visited); err != nil {
			return err
		}

		node := &NodeInstance{
			NodeType: EventNodeType,
			value: EventNode{
				name:        blueprint.Name,
				events:      make([]*com.Event, 0),
				constraints: c,
				okTo:        schedule.nodes[blueprint.OkTo],
				errorTo:     schedule.nodes[blueprint.ErrorTo],
			},
		}

		schedule.nodes[node.value.Name()] = node
		return nil
	}
}

func (schedule *Schedule) addHandlerBlueprint(blueprint com.HandlerBlueprint, visited map[string]bool) error {

	if visited[blueprint.Name] {
		return errors.New("Possible cycle, already visited " + blueprint.Name)
	}
	visited[blueprint.Name] = true

	if _, found := schedule.nodes[blueprint.To]; !found {
		okToEvent, isEvent := schedule.blueprintMaps.Events[blueprint.To]

		if isEvent {
			if err := schedule.addEventBlueprint(okToEvent, visited); err != nil {
				return err
			}
		}

		okToHandler, isHandler := schedule.blueprintMaps.Handlers[blueprint.To]
		if isHandler {
			if err := schedule.addHandlerBlueprint(okToHandler, visited); err != nil {
				return err
			}
		} else {
			// at this point it wasn't found, and it wasn't an event and it wasn't a handler
			return errors.New("Couldn't find the ok-to node for " + blueprint.Name)
		}
	}

	if blueprint.Email.EmailTo != "" {
		node := &NodeInstance{
			NodeType: HandlerNodeType,
			value: EmailHandlerNode{
				emailTo: blueprint.Email.EmailTo,
				to:      schedule.nodes[blueprint.To],
			},
		}

		schedule.nodes[node.value.Name()] = node
	} else {
		return errors.New("Handler " + blueprint.Name + " incorrectly defined.")
	}

	return nil
}
