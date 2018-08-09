package schedule
import (

	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/common"
	"os"
	"io/ioutil"
	"encoding/xml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"log"
	"bytes"
	"fmt"
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

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (fd fileDAO) Save(s Schedule) error {
	
	str := s.Name + ".xml"
	f, err := os.Create(fd.Path + "/" +  str)
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

func (db dbDAO) GetByName(name string) ([]byte, error) {
	var s Schedule
	sEvent := ScheduledEvent{}
	sEvent.ScheduleName = name
	s.Name = name
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	common.CheckError(err)


	rows, err := dbb.NamedQuery(`SELECT * FROM schedules WHERE name=:name`, s)
	if rows == nil {
		log.Fatal(err)
	}
	rows2, err := dbb.NamedQuery(`SELECT * FROM schedulevents WHERE schedulename=:schedulename`, sEvent)
	if rows2 == nil {
		log.Fatal(err)
	}

	sEvents := []ScheduledEvent{}
	
 	for rows2.Next() {
        err := rows2.StructScan(&sEvent)
        common.CheckError(err)
		sEvents = append(sEvents,sEvent)
	} 
	
	for rows.Next() {
        err := rows.StructScan(&s)
        common.CheckError(err)
	}
	
	eventsForSchedule := []Event{}
	for _, e := range sEvents {
		oneEvent := Event{
			Name: e.ScheduleName,
			ReceiveBy: e.EReceiveBy,
		}
		eventsForSchedule = append(eventsForSchedule,oneEvent)
	}

	//encode 
	bytes, err := xml.Marshal(eventsForSchedule)
	if err != nil {
		return nil, err
	}
	s.Schedule = bytes
	schedulebytes, err := xml.Marshal(s)
	log.Println("Our schedule ==========================")
	spew.Dump(s)

	log.Println("Our scheduledEvents ===================")
	spew.Dump(sEvents)

	return schedulebytes,nil
}

func (db dbDAO) Save(s Schedule) error {
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
/* 		dbb.MustExec(schema1)
		dbb.MustExec(schema2) */
		tx := dbb.MustBegin()
		tx.NamedExec("INSERT INTO schedules (name, timing) VALUES (:name,:timing)", &s)
		evnts := []Event{}
		buf := bytes.NewBuffer(s.Schedule)
		dec := xml.NewDecoder(buf)
		var o = Event{}
		for dec.Decode(&o) == nil {
			evnts = append(evnts,o)
		}
		
		fmt.Println("Our scheduled events:")
		spew.Dump(evnts)


		for _, e := range evnts {
		tx.NamedExec("INSERT INTO schedulevents (schedulename, ename, ereceiveby) VALUES (:schedulename, :ename,:ereceiveby)", 
		&ScheduledEvent{
			ScheduleName: s.Name,
			EName: e.Name,
			EReceiveBy: e.ReceiveBy,
			//details
			
		})
		}
		tx.Commit()

	return nil
}