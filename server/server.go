package server

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
//	"io/ioutil"
	"log"
	"net/http"
//	"os"
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

	log.Printf("Received the following information in the event handler: %v\n", event)
	w.WriteHeader(http.StatusOK)

}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	sched := common.Schedule{}
	var fd = database.NewScheduleDAO()

	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["name"]
		if ! ok || len(keys[0]) < 1 {
			log.Println("You didn't have a parameter")
		}


		bytes,err := fd.GetByName(string(keys[0]))
		if err != nil {
			log.Println(err)
			return
		}
 	

		w.Header().Set("Content-Type", "application/xml")

		_,err = w.Write(bytes)

	        if err != nil {
             		fmt.Println("Could not send bytes")
              		return

                }

		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == "PUT" {         
	err := xml.NewDecoder(r.Body).Decode(&sched)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < len(sched.Schedule); i++ {
		valid := validateEvent(sched.Schedule[i])
		if err != nil || valid != nil {
			log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	log.Printf("Received the following information in schedule handler: %v\n", sched)
	
	//var fd = database.NewScheduleDAO()
	err = fd.Save(sched)
	if err != nil {
		log.Println(err)
	}
	}
	w.WriteHeader(http.StatusOK)

}

func validateEvent(e common.Event) error {
	if e.Name == "" {
		return errors.New("Name cannot be empty.")
	} else {
		return nil
	}
}

