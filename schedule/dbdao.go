package schedule

import(
	
	"github.com/jmoiron/sqlx"
	"log"
	"bytes"
	"egbitbucket.dtvops.net/deadline/common"
	"encoding/xml"
)

func (db dbDAO) GetByName(name string) ([]byte, error) {
	var s Schedule
	sEvent := ScheduledEvent{}
	sHandler := ScheduledHandler{}
	//stateless event struct
	sEvent.ScheduleName = name
	s.Name = name
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	common.CheckError(err)

	rows, err := dbb.NamedQuery(`SELECT * FROM schedules WHERE name=:name`, s)
	if rows == nil {
		common.CheckError(err)
	}
	rows2, err := dbb.NamedQuery(`SELECT * FROM schedulevents WHERE schedulename=:schedulename`, sEvent)
	if rows2 == nil {
		common.CheckError(err)
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
	
	rows3, err := dbb.NamedQuery(`SELECT * FROM handlers WHERE schedulename=:schedulename`, &sHandler)
	if rows == nil {
		common.CheckError(err)
	}

	for rows3.Next() {
		err := rows.StructScan(&sHandler)
		common.CheckError(err)
		s.Handler = Handler{
			Name: sHandler.Name,
			Address: sHandler.Address,
		}
	}

	eventsForSchedule := []Event{}
	for _, e := range sEvents {
		oneEvent := Event{
			Name: e.ScheduleName,
			ReceiveBy: e.EReceiveBy,
		}
		eventsForSchedule = append(eventsForSchedule,oneEvent)
	}

	bytes, err := xml.Marshal(eventsForSchedule)
	if err != nil {
		return nil, err
	}
	s.Schedule = bytes
	schedulebytes, err := xml.Marshal(s)

	return schedulebytes,nil
}

func (db dbDAO) Save(s *Schedule) error {
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
		//initalize tables function?

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
		//schedulevents -- stateless event information that correlate to schedules
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

	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil{
		common.CheckError(err)
		return []Schedule{},err
	}
/* 	dbb.MustExec(eventSchema)
	dbb.MustExec(handlerSchema)
	dbb.MustExec(scheduleEventSchema)
	dbb.MustExec(scheduleSchema) */
	
	schedulesFromDB := []Schedule{}
	err = dbb.Select(&schedulesFromDB, "SELECT * FROM schedules")
	if err != nil {
		common.CheckError(err)
		return []Schedule{},err
	}


	//get non-live scheduled event information from that table 
	eventFromDB := ScheduledEvent{} 
	handlerFromDB := ScheduledHandler{}
	eventForSchedule := Event{} 
	eventsForSchedule := []Event{}
	schedulesForTable := []Schedule{}
	for _,s := range schedulesFromDB {
		eventFromDB.ScheduleName = s.Name
		rows, err := dbb.NamedQuery(`SELECT * FROM schedulevents WHERE schedulename=:schedulename`, &eventFromDB)
		if rows == nil {
			common.CheckError(err)
		}
		for rows.Next() {
			err := rows.StructScan(&eventFromDB)
			common.CheckError(err)
			eventForSchedule = Event{
				Name: eventFromDB.EName,
				ReceiveBy: eventFromDB.EReceiveBy,
			}
			eventsForSchedule = append(eventsForSchedule,eventForSchedule)
		} 
		bytes, err := xml.Marshal(eventsForSchedule)

		if err != nil {
			common.CheckError(err)
			return nil, err
		}
		
		s.Schedule = bytes
		eventsForSchedule = []Event{}
		handlerFromDB.ScheduleName = s.Name
		rows, err = dbb.NamedQuery(`SELECT * FROM handlers WHERE schedulename=:schedulename`, &handlerFromDB)
		if rows == nil {
			common.CheckError(err)
		}

		for rows.Next() {
			err := rows.StructScan(&handlerFromDB)
			common.CheckError(err)
			s.Handler = Handler{
				Name: handlerFromDB.Name,
				Address: handlerFromDB.Address,
			}
		}

		schedulesForTable = append(schedulesForTable,s)
	}

	return schedulesForTable,nil
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
	
	dbb, err := sqlx.Open("mysql", db.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	tx := dbb.MustBegin()
	_, err = tx.NamedExec("INSERT INTO events (name, receiveat,success,details, islive) VALUES (:name, :receiveat,:success,:details, :islive)", &e)
	common.CheckError(err)
	return nil
}