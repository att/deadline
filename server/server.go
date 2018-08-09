package server

import (
	"bytes"
	"egbitbucket.dtvops.net/deadline/common"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/schedule"
	"encoding/json"
	"encoding/xml"
	"errors"
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
	common.CheckError(err)
	log.Println(msg)
	w.WriteHeader(http.StatusOK)
}

func eventHander(w http.ResponseWriter, r *http.Request) {

	event := schedule.Event{}
	if r.Body == nil {
		log.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&event)

	valid := event.ValidateEvent()
	if err != nil || valid != nil {
		log.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	M.UpdateEvents(&event)
	w.WriteHeader(http.StatusOK)

}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	err := doMethod(r.Method,w,r)
	if  err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

	} else {
	w.WriteHeader(http.StatusOK) 
	}
}


func doMethod(method string, w http.ResponseWriter, r *http.Request) error {

	sched := schedule.Schedule{}

	if r.Body == nil {
		
		return errors.New("No input")
	}
	switch method {
		case "GET":
			return getSchedule(r,w)
		case "PUT":
			return putSchedule(w,r,sched)
	}
	return nil
}
	


func putSchedule(w http.ResponseWriter, r *http.Request, sched schedule.Schedule)  error {
			err := xml.NewDecoder(r.Body).Decode(&sched)
				if err != nil {
					return err
				}
			var f schedule.Event

			buf := bytes.NewBuffer(sched.Schedule)
			dec := xml.NewDecoder(buf)
			for dec.Decode(&f) == nil {
				e := f
				valid := e.ValidateEvent()
					if valid != nil {
						return valid
					}
				node1 := schedule.Node{
					Event: &e,
					Nodes: []schedule.Node{},
				}
				sched.Start.Nodes = append(sched.Start.Nodes, node1)
			}

			err = Fd.Save(sched)
			common.CheckError(err)
			M.UpdateSchedule(&sched)
			return nil
}

func getSchedule(r *http.Request,w http.ResponseWriter) error {
		keys, ok := r.URL.Query()["name"]
		if !ok || len(keys[0]) < 1 {
			return errors.New("You didn't have a parameter")
		}

		bytes, err := Fd.GetByName(string(keys[0]))
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/xml")

		_, err = w.Write(bytes)

		if err != nil {
			return err
		}

		return nil
		
	}

