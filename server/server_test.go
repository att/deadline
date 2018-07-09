package server

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server = NewDeadlineServer()
var baseAddress = "http://localhost:8081"
var eventApi = baseAddress + "/api/v1/event"

func TestMain(m *testing.M) {
	go server.Start()
	exitVal := m.Run()
	go server.Stop()
	os.Exit(exitVal)
}

func TestGoodParams(test *testing.T) {

	goodRequest := "{\"name\": \"kaela\", \"success\": true}"
	response, err := http.Post(eventApi, "application/json", strings.NewReader(goodRequest))

	assert.Nil(test, err, "Error contacting server")
	assert.NotNil(test, response, "Response object is nil")
	assert.Equal(test, http.StatusOK, response.StatusCode, "Response http status code not what it should be")

}

func TestBadParams(test *testing.T) {

	badReqeust := "{}"
	response, err := http.Post(eventApi, "application/json", strings.NewReader(badReqeust))

	assert.Nil(test, err, "Error contacting server")
	assert.NotNil(test, response, "Response object is nil")
	assert.Equal(test, http.StatusBadRequest, response.StatusCode, "Response http status code not what it should be")

}

func TestGoodEvent(test *testing.T) {
//give good and bad xml files 

}

func TestBadEvent(test *testing.T) {




}
