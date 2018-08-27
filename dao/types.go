package dao

import (
	"github.com/att/deadline/common"
	"github.com/att/deadline/config"
)

func NewScheduleDAO(cfg *config.Config) ScheduleDAO { //TODO fix interface to throw errors
	// if cfg.DAO == "file" {
	return newFileDAO(cfg.FileConfig.Directory)
	// }
	// return &dbDAO{
	// 	ConnectionString: cfg.DBConfig.ConnectionString,
	// }
}

type ScheduleDAO interface {
	GetByName(string) (*ScheduleBlueprint, error)
	Save(s *ScheduleBlueprint) error
	LoadScheduleBlueprints() ([]ScheduleBlueprint, error)
	LoadEvents() ([]common.Event, error)
	SaveEvent(e *common.Event) error
}

type ScheduleBlueprint struct {
	Timing   string             `xml:"timing,attr" db:"timing"`
	Name     string             `xml:"name,attr" db:"name"`
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
	Name        string           `xml:"name,attr"`
	Constraints EventConstraints `xml:"constraints"`
}

type EventConstraints struct {
	ReceiveBy string `xml:"receive-by,omitempty"`
}

type HandlerBlueprint struct {
	Email EmailHandlerBlueprint `xml:"email,omitempty"`
}

type EmailHandlerBlueprint struct {
	To string `xml:"to"`
}
