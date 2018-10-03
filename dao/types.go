package dao

import (
	"time"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

// NewScheduleDAO returns a new NewScheduleDAO based on the given configuration.
func NewScheduleDAO(cfg *config.Config) (ScheduleDAO, error) {
	log = cfg.GetLogger("dao")

	if cfg.Storage == config.FileStorage {
		return newFileDAO(cfg.FileConfig.Directory)
	} else if cfg.Storage == config.DBStorage {
		//
	}

	return newFileDAO(cfg.FileConfig.Directory)
}

// ScheduleDAO is the interface to store and retrieve schedueles from some type of storage.
type ScheduleDAO interface {
	GetByName(string) (*com.ScheduleBlueprint, error)
	Save(s *com.ScheduleBlueprint) error
	LoadScheduleBlueprints() ([]com.ScheduleBlueprint, error)
	EventsAfter(t time.Time) (chan com.Event, error)
	SaveEvent(e *com.Event) error
}
