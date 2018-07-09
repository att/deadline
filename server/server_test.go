package server

import (
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"testing"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"egbitbucket.dtvops.net/deadline/common"
)

var server = NewDeadlineServer()
var baseAddress = "http://localhost:8081"
var eventApi = baseAddress + "/api/v1/event"
var badfile = "badfile.xml"
var goodfile = "sample_schedule.xml"
var testschedule common.Schedule 

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
	xfile, err := os.Open(goodfile)
	assert.Nil(test,err, "Error opening file")
	bytes, err := ioutil.ReadAll(xfile)
	assert.Nil(test, err, "Error getting bytes.")
	assert.NotNil(test, xfile, "XML returns nil")
	err = xml.Unmarshal(bytes, &testschedule)
	assert.Nil(test, err,"Could not decode bytes.")
	
}

func TestBadEvent(test *testing.T) {

        xfile, err := os.Open(badfile)
        assert.Nil(test,err, "Error opening file")
        bytes, err := ioutil.ReadAll(xfile)
        assert.Nil(test, err, "Error getting bytes.")
        assert.NotNil(test, xfile, "XML returns nil")
        err = xml.Unmarshal(bytes, &testschedule)
        assert.NotNil(test, err,"Could not decode bytes.")



}
