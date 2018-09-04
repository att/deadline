package dao

import (
	com "github.com/att/deadline/common"
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
	GetByName(string) (*com.ScheduleBlueprint, error)
	Save(s *com.ScheduleBlueprint) error
	LoadScheduleBlueprints() ([]com.ScheduleBlueprint, error)
	LoadEvents() ([]com.Event, error)
	SaveEvent(e *com.Event) error
}
