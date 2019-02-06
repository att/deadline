package dao

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"encoding/xml"
	"io/ioutil"
	"os"

	com "github.com/att/deadline/common"
)

type fileDAO struct {
	path string
}

func newFileDAO(path string) (*fileDAO, error) {
	dao := &fileDAO{
		path: path,
	}

	// init the directory
	if _, err := makeOrOpenDirectory(dao.path); err != nil {
		return nil, err
	}

	return dao, nil
}

func (dao fileDAO) GetByName(name string) (*com.ScheduleBlueprint, error) {
	var bytes []byte

	file, err := os.Open(dao.path + "/" + name + ".xml")
	defer file.Close()
	if err != nil {
		return nil, err
	}

	if bytes, err = ioutil.ReadAll(file); err != nil {
		return nil, err
	}

	blueprint := &com.ScheduleBlueprint{}
	if xml.Unmarshal(bytes, blueprint); err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (dao fileDAO) Save(blueprint *com.ScheduleBlueprint) error {

	fileName := blueprint.Name + ".xml"
	file, err := os.Create(dao.path + "/" + fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(file)
	if err = encoder.Encode(blueprint); err != nil {
		return err
	}

	return nil

}

func (dao fileDAO) LoadScheduleBlueprints() ([]com.ScheduleBlueprint, error) {
	blueprints := []com.ScheduleBlueprint{}

	directory, err := os.Open(dao.path)
	defer directory.Close()

	if err != nil {
		return blueprints, err
	}

	list, _ := directory.Readdirnames(0)
	for _, scheduleFile := range list {
		if strings.Contains(scheduleFile, ".xml") {
			scheduleName := strings.TrimSuffix(scheduleFile, ".xml")

			if blueprint, err := dao.GetByName(scheduleName); err != nil {
			} else {
				blueprints = append(blueprints, *blueprint)
			}

		}
	}
	return blueprints, nil
}

func (dao fileDAO) EventsAfter(t time.Time) (chan com.Event, error) {
	events := make(chan com.Event)
	after := t.Unix()

	eventDirectory := dao.path + "/" + "events"
	dir, err := makeOrOpenDirectory(eventDirectory)

	if err != nil {
		close(events)
		dir.Close()
		return events, nil
	}

	go func() {

		list, err := dir.Readdirnames(0)
		defer dir.Close()

		if err != nil {
			log.WithError(err).Info("could not read from " + eventDirectory)
			return
		}

		for _, eventFile := range list {
			data, err := ioutil.ReadFile(eventDirectory + "/" + eventFile)
			event := &com.Event{}

			if err == nil {
				if err = json.Unmarshal(data, event); err == nil {
					if event.ReceivedAt >= after {
						events <- *event
					}
				} else {
					log.WithFields(logrus.Fields{
						"error": err,
						"file":  eventFile,
					}).Debug("didn't load event file")
				}
			} else {
				log.WithFields(logrus.Fields{
					"error": err,
					"file":  eventFile,
				}).Debug("didn't load event file")
			}
		}

		close(events)
	}()

	return events, nil
}

func makeOrOpenDirectory(path string) (*os.File, error) {

	if info, err := os.Stat(path); os.IsNotExist(err) { //doesn't exist so make it
		err = os.MkdirAll(path, 0755)

		if err != nil { // couldn't make directory
			return nil, err
		}
	} else if err != nil {
		return nil, err

	} else if !info.IsDir() {
		return nil, errors.New("path " + path + " exists but is not a directory.")
	}

	return os.Open(path)
}

func (dao fileDAO) SaveEvent(e *com.Event) error {
	fileName := e.Name + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".json"

	file, err := os.Create(dao.path + "/events/" + fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(e); err != nil {
		return err
	}

	return nil
}
