package dao

import (
	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
)

// NewScheduleDAO returns a new NewScheduleDAO based on the given configuration.
func NewScheduleDAO(cfg *config.Config) ScheduleDAO { //TODO fix interface to throw errors
	// if cfg.DAO == "file" {
	return newFileDAO(cfg.FileConfig.Directory)
	// }
	// return &dbDAO{
	// 	ConnectionString: cfg.DBConfig.ConnectionString,
	// }
}

// ScheduleDAO is the interface to store and retrieve schedueles from some type of storage.
type ScheduleDAO interface {
	GetByName(string) (*com.ScheduleBlueprint, error)
	Save(s *com.ScheduleBlueprint) error
	LoadScheduleBlueprints() ([]com.ScheduleBlueprint, error)
	LoadEvents() ([]com.Event, error)
	SaveEvent(e *com.Event) error
}
