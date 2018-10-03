package server

import (
	"time"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	com "github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/schedule"
)

var manager *schedule.ScheduleManager
var log *logrus.Logger

// DeadlineServer is the http server for the deadline application.
type DeadlineServer struct {
	server *http.Server
}

// NewDeadlineServer returns a new deadline server based on the configuration given.
func NewDeadlineServer(cfg *config.Config) *DeadlineServer {
	manager = schedule.GetManagerInstance(cfg)
	log = cfg.GetLogger("server")

	return &DeadlineServer{
		server: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: newDeadlineHandler(),
		},
	}
}

// Start starts the http server.
func (dlsvr *DeadlineServer) Start() error {
	return dlsvr.server.ListenAndServe()
}

// Stop stops the http server.
func (dlsvr *DeadlineServer) Stop() error {
	return dlsvr.server.Close()
}

func newDeadlineHandler() http.Handler {

	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/event", eventHandler)
	handler.HandleFunc("/api/v1/blueprint", blueprintHandler)
	handler.HandleFunc("/api/v1/schedule", scheduleHandler)
	return handler
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	event := com.Event{}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&event)

	valid := event.ValidateEvent()
	if err != nil || valid != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event.ReceivedAt = time.Now().Unix()
	manager.Update(&event)
	w.WriteHeader(http.StatusOK)

}

func blueprintHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getBlueprint(w, r)
	case http.MethodPut:
		putBlueprint(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name, err := getParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	schedule := manager.GetSchedule(name)
	if schedule == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(schedule)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)

	if err != nil {
		return
	}
}

func putBlueprint(w http.ResponseWriter, r *http.Request) error {
	blueprint := com.ScheduleBlueprint{}
	err := xml.NewDecoder(r.Body).Decode(&blueprint)

	if err != nil {
		return err
	}

	manager.AddScheduleAndSave(&blueprint)
	w.WriteHeader(http.StatusCreated)

	return nil
}

func getBlueprint(w http.ResponseWriter, r *http.Request) error {
	name, err := getParams(r)
	if err != nil {
		return err
	}

	blueprint, err := manager.GetBlueprint(name) //TODO pull from schedule manager not DAO
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/xml")

	data, err := xml.MarshalIndent(blueprint, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func getParams(r *http.Request) (string, error) {
	keys, ok := r.URL.Query()["name"]
	if !ok || len(keys[0]) < 1 {
		return "", errors.New("You didn't have a parameter")
	}
	return string(keys[0]), nil
}
