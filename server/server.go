package server

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"

	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/schedule"
	
)

var m = schedule.NewManager()
var fd = schedule.NewScheduleDAO()

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

	//log.Printf("Received the following information in the event handler: %v\n", event)
	schedule.UpdateEvents(m, &event, fd)
	w.WriteHeader(http.StatusOK)

}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	sched := schedule.Schedule{}

	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["name"]
		if !ok || len(keys[0]) < 1 {
			log.Println("You didn't have a parameter")
		}

		bytes, err := fd.GetByName(string(keys[0]))
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/xml")

		_, err = w.Write(bytes)

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
		var f common.Event

		buf := bytes.NewBuffer(sched.Schedule)
		dec := xml.NewDecoder(buf)

		for dec.Decode(&f) == nil {
			e := f
			valid := validateEvent(e)
			fmt.Println("Address of an event:")
			fmt.Printf("%p\n", &e)
			if err != nil || valid != nil {
				log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//until we hit eof

			node1 := schedule.Node{
				Event: &e,
				Nodes: []schedule.Node{},
			}
			sched.Start.Nodes = append(sched.Start.Nodes, node1)

		}
		//log.Println("Received the following information in schedule handler: \n", sched)
		fmt.Println("Then we have the nodes that are connected to start")
		fmt.Println(sched.Start.Nodes)
		err = fd.Save(sched)
		if err != nil {
			log.Println(err)
		}
		schedule.UpdateSchedule(m, &sched)
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
