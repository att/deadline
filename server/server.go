package server

import (
	"time"
	"os"
	"bytes"

	"github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/schedule"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

var M *schedule.ScheduleManager


type DeadlineServer struct {
	server *http.Server
}

func NewDeadlineServer(c *config.Config) *DeadlineServer {
	common.Init(os.Stdout,os.Stdout)
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
	handler.HandleFunc("/api/v1/event", eventHandler)
	handler.HandleFunc("/api/v1/schedule", scheduleHandler)
	handler.HandleFunc("/api/v1/msg", notifyHandler)
	handler.HandleFunc("/status", statusHandler)
	handler.HandleFunc("/app/", appHandler)
	return handler
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/app/",http.FileServer(http.Dir("app/src"))).ServeHTTP(w,r)
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	msg := ""
	if r.Body == nil {
		common.Info.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&msg)
	common.CheckError(err)
	common.Info.Println(msg)
	w.WriteHeader(http.StatusOK)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	event := schedule.Event{}
	if r.Body == nil {
		common.Info.Println("No request body sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&event)

	valid := event.ValidateEvent()
	if err != nil || valid != nil {
		common.Info.Println("Cannot accept request. decoding error:", err, "validation error:", valid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event.ReceiveAt = time.Now().Format("15:04:05")
	M.UpdateEvents(&event)
	err = schedule.Fd.SaveEvent(&event)
	common.CheckError(err)
	w.WriteHeader(http.StatusOK)

}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {

	err := doMethod(r.Method,w,r)
	if  err != nil {
		common.Info.Println(err)
		w.WriteHeader(http.StatusBadRequest)

	} else {
	w.WriteHeader(http.StatusOK) 
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		name, err := getParams(r)
		if err != nil {
			common.Info.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		bytes, err := M.GetLiveSchedule(name)
		if err != nil {
			common.Info.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(bytes)

		if err != nil {
			common.Info.Println(err)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
}



func doMethod(method string, w http.ResponseWriter, r *http.Request) error {

	sched := schedule.Definition{}

	if r.Body == nil {
		
		return errors.New("No input")
	}
	switch method {
		case "GET":
			return getSchedule(w,r)
		case "PUT":
			return putSchedule(w,r,sched)
	}
	return nil
}

func putSchedule(w http.ResponseWriter, r *http.Request, sched schedule.Definition)  error {
			err := xml.NewDecoder(r.Body).Decode(&sched)
				if err != nil {
					return err
				}
			var evnt schedule.Event

			buf := bytes.NewBuffer(sched.ScheduleContent)
			dec := xml.NewDecoder(buf)
			for dec.Decode(&evnt) == nil {
				e := evnt
				//change location memory
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
			M.UpdateSchedule(&sched)
			return nil
}

func getSchedule(w http.ResponseWriter,r *http.Request) error {
		name, err := getParams(r)
		if err != nil {
			return err 
		}

		bytes, err := schedule.Fd.GetByName(name)
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


func getParams(r *http.Request) (string, error){
	keys, ok := r.URL.Query()["name"]
	if !ok || len(keys[0]) < 1 {
		return "",errors.New("You didn't have a parameter")
	}
	return string(keys[0]), nil
}

