package server

import (
	"bytes"
	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/schedule"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var M *schedule.ScheduleManager
var Fd schedule.ScheduleDAO

type DeadlineServer struct {
	server *http.Server
}

func NewDeadlineServer(c *config.Config) *DeadlineServer {

	return &DeadlineServer{
		server: &http.Server{
			Addr:    ":" + c.Server.Port,
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
	handler.HandleFunc("/api/v1/msg", notifyHandler)
	return handler
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	msg := ""
	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Println(err)
	}
	log.Println(msg)
	w.WriteHeader(http.StatusOK)
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
	schedule.UpdateEvents(M, &event, Fd)
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

		bytes, err := Fd.GetByName(string(keys[0]))
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

			if err != nil || valid != nil {
				log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			node1 := schedule.Node{
				Event: &e,
				Nodes: []schedule.Node{},
			}
			sched.Start.Nodes = append(sched.Start.Nodes, node1)

		}

		err = Fd.Save(sched)
		if err != nil {
			log.Println(err)
		}
		schedule.UpdateSchedule(M, &sched)
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
