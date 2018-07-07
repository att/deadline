package server

import (
	"encoding/json"
	"log"
	"net/http"

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
		http.Error(w, "", http.StatusBadRequest)
	}

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Println("Could not decode", err)
		http.Error(w, "", http.StatusBadRequest)
	}

	log.Printf("Received the following information: %v\n", event)

	w.WriteHeader(http.StatusOK)

}
