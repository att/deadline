package database

import (
	//"github.com/mattn/go-sqlite3"
	//"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"egbitbucket.dtvops.net/deadline/common"
	//"github.com/Masterminds/structable"
)

//possibly use structable?

type ScheduleDAO interface {
	getByName() (common.Schedule, error)
	save(s common.Schedule) error
}

type fileDAO struct {
}

func (fd fileDAO) getByName() (common.Schedule, error) {

	var s common.Schedule
	//get a name from a directory

	o, err := os.Open(s.Name + ".xml")
	if err != nil {
		return common.Schedule{}, err
	}

	defer o.Close()

	//read in the xml file
	bytes, _ := ioutil.ReadAll(o)
	err = xml.Unmarshal(bytes, &s)

	if err != nil {
		return common.Schedule{}, err
	}

	fmt.Printf("We have the following struct: %#v\n", s)
	return s, nil

}

func (fd fileDAO) save(s common.Schedule) error {

	f, err := os.Create(s.Name + ".xml")
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

func NewSchedule() ScheduleDAO {
	// return some new object
	//if wwe have the indicator we are using a datbase, return database struct
	//eeelse we will use files
	return &fileDAO{}
}

//type dbDAO struct {}

/*

func (d dbDAO) getByName common.Schedule, error {


//get a row, return a struct
var s common.Schedule
//get a name from a directory


err := db.QueryRow("SELECT * from s")
for rows.Next() {
        err := rows.StructScan(&s)
        if err != nil {return &common.Schedule{}, err}
	}
}

func (d dbDAO) save(s common.Schedule) error {

//insert/update a row based on the schedule given
//r := structable.New(db, "sqlite3").Bind("test_table", s)
//err := r.Insert()



}
*/
