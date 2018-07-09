package server

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
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

	return handler
}

func eventHander(w http.ResponseWriter, r *http.Request) {

	event := common.Event{}

	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&event)
	valid := validateEvent(event)
	if err != nil || valid != nil {
		log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Received the following information: %v\n", event)

	w.WriteHeader(http.StatusOK)

}

func validateEvent(e common.Event) error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

func getEvent (s common.Schedule)  error {

	file, err := os.Open("sample_schedule.xml")
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
	
	for i := 0; i < len(s.Schedule); i++ {

	fmt.Println (s.Schedule[i].Name)
	}

	return nil

}


