package common

import "errors"

type ScheduleBlueprint struct {
	Timing   string             `xml:"timing,attr" db:"timing"`
	Name     string             `xml:"name,attr" db:"name"`
	StartsAt string             `xml:"starts-at,attr" db:"name"`
	Events   []EventBlueprint   `xml:"event"`
	Handlers []HandlerBlueprint `xml:"handler"`
	Start    StartBlueprint     `xml:"start"`
	End      EndBlueprint       `xml:"end"`
}

type StartBlueprint struct {
	To string `xml:"to,attr"`
}

type EndBlueprint struct {
	Name string `xml:"name,attr"`
}

type EventBlueprint struct {
	Name        string                    `xml:"name,attr"`
	Constraints EventConstraintsBlueprint `xml:"constraints"`
	OkTo        string                    `xml:"ok,attr"`
	ErrorTo     string                    `xml:"error,attr"`
}

type EventConstraintsBlueprint struct {
	ReceiveBy string `xml:"receive-by,omitempty"`
}

type HandlerBlueprint struct {
	Name  string                `xml:"name,attr"`
	Email EmailHandlerBlueprint `xml:"email,omitempty"`
	To    string                `xml:"to,attr"`
}

type EmailHandlerBlueprint struct {
	EmailTo string `xml:"to"`
}

type BlueprintMaps struct {
	Events   map[string]EventBlueprint
	Handlers map[string]HandlerBlueprint
}

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
