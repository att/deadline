package schedule

import(
	"log"
	"bytes"
	"time"
	"egbitbucket.dtvops.net/deadline/common"
	"github.com/jmoiron/sqlx"
	"encoding/xml"
	
)

func (db dbDAO) GetByName(name string) ([]byte, error) {
	var s Schedule
	s.Name = name
	sEvent := ScheduledEvent{}
	sEvent.ScheduleName = s.Name
	sHandler := ScheduledHandler{}
	sHandler.ScheduleName = s.Name

	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	common.CheckError(err)

	rows, err := dbb.NamedQuery(`SELECT * FROM schedules WHERE name=:name`, s)
	common.CheckError(err)

	for rows.Next() {
        err := rows.StructScan(&s)
        common.CheckError(err)
	}

	rows2, err := dbb.NamedQuery(`SELECT * FROM schedulevents WHERE schedulename=:schedulename`, sEvent)
	common.CheckError(err)

	eventsForSchedule := []Event{}
	for rows2.Next() {
		err := rows2.StructScan(&sEvent)
		common.CheckError(err)
		eventForSchedule := Event{
			Name: sEvent.EName,
			ReceiveBy: sEvent.EReceiveBy,
		}
		eventsForSchedule = append(eventsForSchedule,eventForSchedule)
	} 

	bytes, err := xml.Marshal(eventsForSchedule)
	common.CheckError(err)
	s.Schedule = bytes
	
	rows3, err := dbb.NamedQuery(`SELECT * FROM handlers WHERE schedulename=:schedulename`, &sHandler)
	common.CheckError(err)
	
	for rows3.Next() {
		err := rows3.StructScan(&sHandler)
		common.CheckError(err)

		s.Handler = Handler{
			Name: sHandler.Name,
			Address: sHandler.Address,
		}
	}

	schedulebytes, err := xml.Marshal(s)
	return schedulebytes, err
}

func (db dbDAO) Save(s *Schedule) error {
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

		tx := dbb.MustBegin()
		_, err = tx.NamedExec("INSERT INTO schedules (name, timing) VALUES (:name,:timing)", &s)
		common.CheckError(err)
		evnts := []Event{}
		buf := bytes.NewBuffer(s.Schedule)
		dec := xml.NewDecoder(buf)
		var o = Event{}
		for dec.Decode(&o) == nil {
			evnts = append(evnts,o)
		}
		for _, e := range evnts {
		_, err = tx.NamedExec("INSERT INTO schedulevents (schedulename, ename, ereceiveby) VALUES (:schedulename, :ename,:ereceiveby)", 
		&ScheduledEvent{
			ScheduleName: s.Name,
			EName: e.Name,
			EReceiveBy: e.ReceiveBy,
		})
		common.CheckError(err)
		}

		scheduleHandler := s.Handler
		handlerForDB := ScheduledHandler{
			ScheduleName: s.Name,
			Name:	scheduleHandler.Name,
			Address: scheduleHandler.Address,

		}
		_, err = tx.NamedExec("INSERT INTO handlers (schedulename, name, address) VALUES (:schedulename, :name,:address)", &handlerForDB)
		common.CheckError(err)
		

		tx.Commit()

	return nil
}

func (db dbDAO) LoadStatelessSchedules() ([]Schedule,error){
	schedulesFromDB := []Schedule{}

	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil{
		common.CheckError(err)
		return []Schedule{},err
	}

	//initializeTables(dbb)

	err = dbb.Select(&schedulesFromDB, "SELECT * FROM schedules")
	if err != nil {
		common.CheckError(err)
		return []Schedule{},err
	}

	for s := 0; s < len(schedulesFromDB); s++ {
		bytes, err := db.GetByName(schedulesFromDB[s].Name)
		common.CheckError(err)
		err = xml.Unmarshal(bytes,&schedulesFromDB[s])
		common.CheckError(err)
	}

	return schedulesFromDB,nil
}



func (db dbDAO) LoadEvents() ([]Event,error){
	
	liveEvents := []Event{}
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil{
		return []Event{},err
	}
	err = dbb.Select(&liveEvents,`SELECT * FROM events`)
	if err != nil {
		common.CheckError(err)
		return []Event{}, err
	}
	return liveEvents, nil
}

func (db dbDAO) SaveEvent(e *Event) error{
	
	e.ReceiveAt = time.Now().Format("01-02-2006 15:04:05")
	e.IsLive = true
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	tx := dbb.MustBegin()
	_, err = tx.NamedExec("INSERT INTO events (name, receiveat,success,islive) VALUES (:name, :receiveat,:success,:islive)", e)
	common.CheckError(err)
	tx.Commit()
	return nil
}

func initializeTables(dbb *sqlx.DB) {
	dbb.MustExec(eventSchema)
	dbb.MustExec(handlerSchema)
	dbb.MustExec(scheduleEventSchema)
	dbb.MustExec(scheduleSchema) 
}
