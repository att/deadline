package server

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"egbitbucket.dtvops.net/deadline/schedule"
	"egbitbucket.dtvops.net/deadline/common"
)

type DeadlineServer struct {
	server *http.Server
}

func NewDeadlineServer() *DeadlineServer {
	return &DeadlineServer{
		server: &http.Server{
			Addr:    ":8081",
			Handler: newDeadlineHandler(),
		},
	}
}

func (dlsvr *DeadlineServer) Start() error {
	return dlsvr.server.ListenAndServe()
}

func (dlsvr *DeadlineServer) Stop() error {
	return dlsvr.server.Close()
}

func newDeadlineHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/event", eventHander)
	handler.HandleFunc("/api/v1/schedule", scheduleHandler)
	return handler
}

func eventHander(w http.ResponseWriter, r *http.Request) {

	event := common.Event{}
	//sched := common.Schedule{}
	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&event)
	//err2 := xml.NewDecoder(r.Body).Decode(&sched)

	valid := validateEvent(event)
	if err != nil || valid != nil {
		log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Received the following information in the event handler: %v\n", event)
	//log.Printf("%v\n",sched)
	w.WriteHeader(http.StatusOK)

}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	sched := common.Schedule{}

	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := xml.NewDecoder(r.Body).Decode(&sched)
	if err != nil {w.WriteHeader(http.StatusBadRequest)
	return}
	for i := 0; i < len(sched.Schedule); i++ {
		valid := validateEvent(sched.Schedule[i])
		if err != nil || valid != nil {
			log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	log.Printf("Received the following information in schedule handler: %v\n", sched)
	var fd = database.NewScheduleDAO()
	fd.Save(sched)
	
	w.WriteHeader(http.StatusOK)

}

func validateEvent(e common.Event) error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

func GetEvent(s *common.Schedule) error {

	file, err := os.Open("sampe_schedule.xml")
	if err != nil {
		return errors.New("Could not open file.")
	}

	defer file.Close()

	//read in the xml file
	bytes, _ := ioutil.ReadAll(file)
	err = xml.Unmarshal(bytes, &s)
	if err != nil {
		return errors.New("Could not make struct.")
	}

	//then we will print all of the events in the schedule
	fmt.Println("Looking at our schedule: ")
	//prints out the schedule of events (names only)

	for i := 0; i < len(s.Schedule); i++ {
		fmt.Println(i)
		fmt.Println(s.Schedule[i].Name)
	}

	return nil

}
