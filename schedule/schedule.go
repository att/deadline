package schedule

import (
	"errors"
	"sync"
	"time"

	com "github.com/att/deadline/common"
)

// Evaluate the schedule completely.
func (schedule *Schedule) Evaluate() State {
	schedule.walk(schedule.Start)

	return schedule.state
}

func (schedule *Schedule) walk(instance *NodeInstance) {
	if schedule.state != Running {
		return
	}

	switch node := instance.value.(type) {

	case *StartNode:
		schedule.walk(node.to)
	case EventNode:
		next, _ := node.Next()
		if next[0] == node.errorTo {
			schedule.state = Failed
		}

		schedule.walk(next[0])

	case *EmailHandlerNode:
		//handle
		schedule.walk(node.to)
	case *EndNode:
		if schedule.state != Failed {
			schedule.state = Ended
		}
	case nil:
		log.Info("nil node type")
	default:
		log.Info("unknown node type")
	}
}

// EventOccurred is the interface into the schedule to tell it that an event
// has occured. The schedule must find the appropriate node that's listening for
// the event and update it.
func (schedule *Schedule) EventOccurred(event com.Event) {

	for _, node := range schedule.nodes {
		if node.value.Name() == event.Name {

			if eventNode, ok := node.value.(EventNode); ok {
				schedule.eventLock.Lock()
				eventNode.AddEvent(event)
				schedule.eventLock.Unlock()

			} else {
				// log cast failure
			}

		}
	}

}

// SubscribesTo returns the set of named events that this schedule subscribes to.  That is,
// all the event nodes that this schedule expects to accept
func (schedule *Schedule) SubscribesTo() map[string]bool {
	return schedule.subscribesTo
}

// FromBlueprint creates a Schedule object from a blueprint.  Errors can occur for various reasons
// like invalid business rules like an event node's error-to can only go to end or a handler
// or just general malformed blueprints like nodes having cycles or hanging nodes.
func FromBlueprint(blueprint *com.ScheduleBlueprint) (*Schedule, error) {
	maps := &com.BlueprintMaps{}
	var err error
	var startTime time.Time

	if maps, err = com.GetBlueprintMaps(blueprint); err != nil {
		return nil, err
	} else if err = checkEmptyFields(blueprint); err != nil {
		return nil, err
	} else if startTime, err = time.Parse(time.RFC3339, blueprint.StartsAt); err != nil {
		return nil, err
	}

	schedule := &Schedule{
		nodes:         make(map[string]*NodeInstance),
		subscribesTo:  make(map[string]bool),
		blueprintMaps: *maps,
		Name:          blueprint.Name,
		StartTime:     startTime,
		eventLock:     &sync.RWMutex{},
		state:         Running,
	}

	if blueprint.End.Name == "" {
		return nil, errors.New("end node must be specified with a valid name")
	}

	schedule.End = &NodeInstance{
		NodeType: EndNodeType,
		value: &EndNode{
			name: blueprint.End.Name,
		},
	}

	schedule.nodes[blueprint.End.Name] = schedule.End
	visited := make(map[string]bool)
	var firstEvent *NodeInstance
	var found bool

	if firstEvent, found := maps.Events[blueprint.Start.To]; !found {
		return nil, errors.New("start node needs to point to an event Node")
	} else if err := schedule.addEventBlueprint(firstEvent, visited); err != nil {
		return nil, err
	}

	if firstEvent, found = schedule.nodes[blueprint.Start.To]; !found {
		return nil, errors.New("schedule built, but still no first node")
	}

	startNode := &NodeInstance{
		NodeType: StartNodeType,
		value: &StartNode{
			to: firstEvent,
		},
	}

	schedule.nodes[startNode.value.Name()] = startNode
	schedule.Start = startNode

	if err := schedule.createdAll(); err != nil {
		return nil, err
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
			// found it, but it's the end node. that's ok
		} else {
			return errors.New("events can only ok-to other events or the end node")
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
			return errors.New("couldn't find the error-to node for " + blueprint.Name)
		}
	}

	return nil
}

// helper function to add an event blueprint to the schedule. This function may recursively call itself while updating
// the visited map to indicate that the current node has been visited.  Can throw an error for various reasons.
func (schedule *Schedule) addEventBlueprint(blueprint com.EventBlueprint, visited map[string]bool) error {
	var c com.EventConstraints
	var err error

	if c, err = com.FromBlueprint(schedule.StartTime, blueprint.Constraints); err != nil {
		return err
	} else if blueprint.Name == "" {
		return errors.New("node names cannot be empty")
	} else if visited[blueprint.Name] {
		return errors.New("possible cycle, already visited " + blueprint.Name)
	}

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
			events:      make([]com.Event, 0),
			constraints: c,
			okTo:        schedule.nodes[blueprint.OkTo],
			errorTo:     schedule.nodes[blueprint.ErrorTo],
		},
	}

	schedule.nodes[node.value.Name()] = node
	schedule.subscribesTo[node.value.Name()] = true
	return nil
}

// helper function to add a handler blueprint to the schedule. This function may recursively call itself while updating
// the visited map to indicate that the current node has been visited.  Can throw an error for various reasons.
func (schedule *Schedule) addHandlerBlueprint(blueprint com.HandlerBlueprint, visited map[string]bool) error {

	if visited[blueprint.Name] {
		return errors.New("possible cycle, already visited " + blueprint.Name)
	} else if blueprint.Name == "" {
		return errors.New("names of nodes cannot be empty")
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
				name:    blueprint.Name,
			},
		}

		schedule.nodes[node.value.Name()] = node
	} else {
		return errors.New("Handler " + blueprint.Name + " incorrectly defined.")
	}

	return nil
}

// helper function to be sure we've created all the nodes requried.
func (schedule *Schedule) createdAll() error {
	for _, event := range schedule.blueprintMaps.Events {
		if _, created := schedule.nodes[event.Name]; !created {
			return errors.New("didnt create event " + event.Name + ", no route to node")
		}
	}

	for _, handler := range schedule.blueprintMaps.Handlers {
		if _, created := schedule.nodes[handler.Name]; !created {
			return errors.New("didnt create handler " + handler.Name + ", no route to node")
		}
	}

	return nil
}

func checkEmptyFields(blueprint *com.ScheduleBlueprint) error {
	if blueprint.Name == "" {
		return errors.New("node names cannot be empty")
	} else if blueprint.StartsAt == "" {
		return errors.New("starts at time cannot be empty")
	} else if blueprint.Timing == "" {
		return errors.New("timing cannot be empty")
	} else {
		return nil
	}
}
