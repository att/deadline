package common

import "errors"

// ScheduleBlueprint is the blueprint struct for a schedule. It is the serialable and static version of a
// schedule.
type ScheduleBlueprint struct {
	Timing   string             `xml:"timing,attr" db:"timing"`
	Name     string             `xml:"name,attr" db:"name"`
	StartsAt string             `xml:"starts-at,attr" db:"name"`
	Events   []EventBlueprint   `xml:"event"`
	Handlers []HandlerBlueprint `xml:"handler"`
	Start    StartBlueprint     `xml:"start"`
	End      EndBlueprint       `xml:"end"`
}

// StartBlueprint is a blueprint for a start node.
type StartBlueprint struct {
	To string `xml:"to,attr"`
}

// EndBlueprint is a blueprint for an end node.
type EndBlueprint struct {
	Name string `xml:"name,attr"`
}

// EventBlueprint is a blueprint for an event node.
type EventBlueprint struct {
	Name        string                    `xml:"name,attr"`
	Constraints EventConstraintsBlueprint `xml:"constraints"`
	OkTo        string                    `xml:"ok,attr"`
	ErrorTo     string                    `xml:"error,attr"`
}

// EventConstraintsBlueprint is a blueprint for an EventConstraints struct.
type EventConstraintsBlueprint struct {
	ReceiveBy string `xml:"receive-by,omitempty"`
}

// HandlerBlueprint is a blueprint for handler node
type HandlerBlueprint struct {
	Name  string                `xml:"name,attr"`
	Email EmailHandlerBlueprint `xml:"email,omitempty"`
	To    string                `xml:"to,attr"`
}

// EmailHandlerBlueprint is a blueprint for an email handler
type EmailHandlerBlueprint struct {
	EmailTo string `xml:"to"`
}

// BlueprintMaps is a helper struct to keep event and handler blueprints in a single struct
type BlueprintMaps struct {
	Events   map[string]EventBlueprint
	Handlers map[string]HandlerBlueprint
}

// GetBlueprintMaps returns BlueprintMaps for a given blueprint for a schedule
func GetBlueprintMaps(blueprint *ScheduleBlueprint) (*BlueprintMaps, error) {

	eventMap := make(map[string]EventBlueprint)
	handlerMap := make(map[string]HandlerBlueprint)

	for _, event := range blueprint.Events {
		if _, found := eventMap[event.Name]; found {
			return nil, errors.New("Two or more nodes use the same name " + event.Name)
		} else {
			eventMap[event.Name] = event
		}
	}

	for _, handler := range blueprint.Handlers {
		if _, found := handlerMap[handler.Name]; found {
			return nil, errors.New("Two or more nodes use the same name " + handler.Name)
		} else {
			handlerMap[handler.Name] = handler
		}
	}

	return &BlueprintMaps{
		Events:   eventMap,
		Handlers: handlerMap,
	}, nil
}
