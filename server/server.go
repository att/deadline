package server

import (
	"net/http"
)

type deadlineServer struct {
	serv1 *http.Server
}

type Event struct {
	Name    string            `json:"name"`
	Success bool              `json:"success"`
	Details map[string]string `json:"details,omitempty"`
}

//log.Fatal(http.ListenAndServe(ds.serv1.Addr, nil))
