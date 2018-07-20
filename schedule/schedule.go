package database

import (
	//"github.com/mattn/go-sqlite3"
	//"database/sql"
	"encoding/xml"
	"fmt"
	//"strings"
	"io/ioutil"
	"os"
	"egbitbucket.dtvops.net/deadline/common"
)

type scheduleManager struct {

	Manager map[string][]common.Schedule

}

func NewManager() *scheduleManager {
	return &scheduleManager{
	Manager: make(map[string][]common.Schedule),

	}

}

type ScheduleDAO interface {
	GetByName(string) ([]byte, error)
	Save(s common.Schedule) error
}

type fileDAO struct {
}

func (fd fileDAO) GetByName(name string) ([]byte, error) {

//	var s common.Schedule
	
                file, err := os.Open(name + ".xml")
                if err != nil {
                        
                        return nil,err

                        }
                fmt.Println("We could open it!")
                defer file.Close()


                //read in the xml file
                bytes, err := ioutil.ReadAll(file)
                if err != nil {
                        
                 	return nil,err

                        }
	return bytes, nil
}

func (fd fileDAO) Save(s common.Schedule) error {
	str := s.Name + ".xml"
	f, err := os.Create(str)
	if err != nil {
		return err
	}
	defer f.Close()
	
	encoder := xml.NewEncoder(f)
	err = encoder.Encode(s)

	if err != nil {
		return err
	}
	return nil
}

func NewScheduleDAO() ScheduleDAO {
	// return some new object
	//if wwe have the indicator we are using a datbase, return database struct
	//eeelse we will use files
	return &fileDAO{}
}

func updateEvents(m *scheduleManager, e common.Event) {
//once you receive an event, tell every schedule that you have it by adding it to their array
var scheds[]common.Schedule = m.Manager[e.Name]
if scheds == nil {return}
    for _, sched := range scheds {
       	sched.ReceivedEvents = append(sched.ReceivedEvents,e)
    }


}

func updateSchedule(m *scheduleManager, s common.Schedule) {
//loop through array and subscribe to every event, and then add itself to the map for every event

        for i := 0; i < len(s.Schedule); i++ {
                //subscribe to every event
                //put into map
                var scheds[]common.Schedule = m.Manager[(s.Schedule[i].Name)]
		if scheds == nil {
		m.Manager[(s.Schedule[i].Name)]=[]common.Schedule{s}
		continue
		}
		scheds = append(scheds,s)
		m.Manager[s.Schedule[i].Name]=scheds

        }
    	fmt.Println("-----------------------------------------------------------") 
//print the map that we have
	fmt.Printf("%#v\n", m.Manager)
	fmt.Println("-----------------------------------------------------------")

}




//type dbDAO struct {}

