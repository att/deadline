package server

import (
	"time"

	"github.com/att/deadline/dao"

	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"

	"github.com/att/deadline/common"
	"github.com/att/deadline/config"
	"github.com/att/deadline/schedule"
)

var M *schedule.ScheduleManager

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
	handler.HandleFunc("/api/v1/event", eventHandler)
	handler.HandleFunc("/api/v1/blueprint", blueprintHandler)
	handler.HandleFunc("/api/v1/schedule", scheduleHandler)
	return handler
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	event := common.Event{}
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
		common.Info.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	schedule := M.GetSchedule(name)
	if schedule == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	bytes, err := json.Marshal(schedule)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)

	if err != nil {
		common.Info.Println(err)
		return
	}
}

func putBlueprint(w http.ResponseWriter, r *http.Request) error {
	blueprint := dao.ScheduleBlueprint{}
	err := xml.NewDecoder(r.Body).Decode(&blueprint)

	if err != nil {
		return err
	}

	M.AddSchedule(&blueprint)
	w.WriteHeader(http.StatusCreated)

	return nil
}

func getBlueprint(w http.ResponseWriter, r *http.Request) error {
	name, err := getParams(r)
	if err != nil {
		return err
	}

	blueprint, err := schedule.Fd.GetByName(name) //TODO pull from schedule manager not DAO
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
