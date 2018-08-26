package schedule
import (
	"strings"

	"github.com/att/deadline/config"
	"github.com/att/deadline/common"
	"os"
	"io/ioutil"
	"encoding/xml"
	_ "github.com/go-sql-driver/mysql"

)

func NewScheduleDAO(c *config.Config) ScheduleDAO {
	if (c.DAO == "file"){
	return &fileDAO{
		Path: c.FileConfig.Directory,
	} 
}
	return &dbDAO{
		ConnectionString: c.DBConfig.ConnectionString,
	}
}


func (fd fileDAO) GetByName(name string) ([]byte, error) {
	file, err := os.Open(fd.Path + "/" + name + ".xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()


	return ioutil.ReadAll(file)
}

func (fd fileDAO) Save(s *Definition) error {
	
	str := s.Name + ".xml"
	f, err := os.Create(fd.Path + "/" +  str)
	if err != nil {
		return err
	}
	defer f.Close()
	s.fixSchedule()
	encoder := xml.NewEncoder(f)
	err = encoder.Encode(s)

	if err != nil {
		return err
	}
	return nil
}

func (fd fileDAO) LoadSchedules() ([]Live,error) { 
	var schedules = []Live{} //live
	s := Definition{}				//live
	file, err := os.Open(fd.Path)
	if err != nil {
		common.Info.Println("Could not open directory.")
		return []Live{}, err
	}
	defer file.Close()

	list,_ := file.Readdirnames(0)
	for _, schedule := range list {
		if strings.Contains(schedule,".xml") {
			schedule = strings.TrimSuffix(schedule,".xml")
			bytes,_ := fd.GetByName(schedule)
			err = xml.Unmarshal(bytes,&s)
			if err != nil {
				return nil, err
			}
			s.MakeNodes()
			l := s.ConvertToLive()
			schedules = append(schedules,*l)
		}
	}
	return schedules,nil
}


func (fd fileDAO) LoadEvents() ([]Event,error){
	liveEvents := []Event{}
	liveEvent := Event{}
	file, err := os.Open(fd.Path + "/" + "events") 
	if err != nil {
		common.Info.Println("Could not open directory.")
		return []Event{}, err
	}
	defer file.Close()

	list,_ := file.Readdirnames(0)
	for _, event := range list {
		if strings.Contains(event,".xml") {
			event = strings.TrimSuffix(event,".xml")
			bytes,_ := fd.GetByName("events/" + event)
			err = xml.Unmarshal(bytes,&liveEvent)
			if err != nil {
				common.Info.Println(event + " wasn't translated")
				continue
			}
			liveEvents = append(liveEvents,liveEvent)
		}
	}


	return liveEvents, nil
}

func (fd fileDAO) SaveEvent(e *Event) error{
	str := e.Name + ".xml"
	f, err := os.Create(fd.Path + "/events/" +  str)
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

