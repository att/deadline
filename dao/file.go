package dao

import (
	"errors"
	"strings"

	"encoding/xml"
	"io/ioutil"
	"os"

	com "github.com/att/deadline/common"
)

type fileDAO struct {
	path string
}

func newFileDAO(path string) *fileDAO {
	dao := &fileDAO{
		path: path,
	}

	// init the directory
	makeOrOpenDirectory(dao.path) // if err, return it

	return dao
}

func (fd fileDAO) GetByName(name string) (*com.ScheduleBlueprint, error) {

	file, err := os.Open(fd.path + "/" + name + ".xml")
	defer file.Close()
	if err != nil {
		return nil, err
	}

	if bytes, err := ioutil.ReadAll(file); err != nil {
		return nil, err
	} else {
		blueprint := &com.ScheduleBlueprint{}
		if xml.Unmarshal(bytes, blueprint); err != nil {
			return nil, err
		} else {
			return blueprint, nil
		}
	}
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
	} else {
		return nil
	}

}

func (dao fileDAO) LoadScheduleBlueprints() ([]com.ScheduleBlueprint, error) {
	blueprints := []com.ScheduleBlueprint{}
	//blueprint := ScheduleBlueprint{}

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

func (fd fileDAO) LoadEvents() ([]com.Event, error) {
	liveEvents := []com.Event{}
	// liveEvent := common.Event{}
	// file, err := makeOrOpenDirectory(fd.path + "/" + "events")
	// defer file.Close()

	// if err != nil {
	// 	common.Info.Println("Cannot read events because", err)
	// 	return []common.Event{}, err
	// }

	// list, _ := file.Readdirnames(0)
	// for _, event := range list {
	// 	if strings.Contains(event, ".xml") {
	// 		event = strings.TrimSuffix(event, ".xml")
	// 		bytes, _ := fd.GetByName("events/" + event)
	// 		err = xml.Unmarshal(bytes, &liveEvent)
	// 		if err != nil {
	// 			common.Info.Println(event + " wasn't translated")
	// 			continue
	// 		}
	// 		liveEvents = append(liveEvents, liveEvent)
	// 	}
	// }

	return liveEvents, nil
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

func (fd fileDAO) SaveEvent(e *com.Event) error {
	str := e.Name + ".xml"
	f, err := os.Create(fd.path + "/events/" + str)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := xml.NewEncoder(f)
	err = encoder.Encode(e)

	if err != nil {
		return err
	}

	return nil
}
