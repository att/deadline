package dao

import (
	"errors"
	"strings"

	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/att/deadline/common"
	_ "github.com/go-sql-driver/mysql"
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

func (fd fileDAO) GetByName(name string) ([]byte, error) {

	file, err := os.Open(fd.path + "/" + name + ".xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (dao fileDAO) Save(s *common.Definition) error {

	str := s.Name + ".xml"
	f, err := os.Create(dao.path + "/" + str)
	defer f.Close()

	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(f)
	err = encoder.Encode(s)

	if err != nil {
		return err
	}
	return nil
}

func (fd fileDAO) LoadSchedules() ([]common.Definition, error) {
	var schedules = []common.Definition{}
	s := common.Definition{}
	file, err := os.Open(fd.path)
	if err != nil {
		common.Info.Println("Could not open directory.")
		return []common.Definition{}, err
	}
	defer file.Close()

	list, _ := file.Readdirnames(0)
	for _, schedule := range list {
		if strings.Contains(schedule, ".xml") {
			schedule = strings.TrimSuffix(schedule, ".xml")
			bytes, _ := fd.GetByName(schedule)
			err = xml.Unmarshal(bytes, &s)
			if err != nil {
				common.Info.Println(schedule + " wasn't translated")
				continue
			}
			schedules = append(schedules, s)
		}
	}
	return schedules, nil
}

func (fd fileDAO) LoadEvents() ([]common.Event, error) {
	liveEvents := []common.Event{}
	liveEvent := common.Event{}
	file, err := makeOrOpenDirectory(fd.path + "/" + "events")
	defer file.Close()

	if err != nil {
		common.Info.Println("Cannot read events because", err)
		return []common.Event{}, err
	}

	list, _ := file.Readdirnames(0)
	for _, event := range list {
		if strings.Contains(event, ".xml") {
			event = strings.TrimSuffix(event, ".xml")
			bytes, _ := fd.GetByName("events/" + event)
			err = xml.Unmarshal(bytes, &liveEvent)
			if err != nil {
				common.Info.Println(event + " wasn't translated")
				continue
			}
			liveEvents = append(liveEvents, liveEvent)
		}
	}

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

func (fd fileDAO) SaveEvent(e *common.Event) error {
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
