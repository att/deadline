package dao

import (
	"github.com/att/deadline/common"
	"github.com/att/deadline/config"
)

func NewScheduleDAO(cfg *config.Config) ScheduleDAO { //TODO fix interface to throw errors
	if cfg.DAO == "file" {
		return newFileDAO(cfg.FileConfig.Directory)
	}
	return &dbDAO{
		ConnectionString: cfg.DBConfig.ConnectionString,
	}
}

type ScheduleDAO interface {
	GetByName(string) ([]byte, error)
	Save(s *common.Definition) error
	LoadSchedules() ([]common.Definition, error)
	LoadEvents() ([]common.Event, error)
	SaveEvent(e *common.Event) error
}

type dbDAO struct {
	ConnectionString string
}
