package server

// import (
// 	"bytes"
// 	"encoding/xml"
// 	"io/ioutil"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"testing"

// 	"github.com/att/deadline/dao"

// 	"github.com/att/deadline/config"
// 	"github.com/stretchr/testify/assert"
// )

// var c = config.Config{
// 	FileConfig: config.FileConfig{
// 		Directory: os.TempDir() + "/deadline_test/" + strconv.Itoa(rand.Int()),
// 	},
// 	DAO: "file",
// 	Server: config.ServConfig{
// 		Port: "8081",
// 	},
// }
// var server = NewDeadlineServer(&c)
// var baseAddress = "http://localhost:8081"
// var eventApi = baseAddress + "/api/v1/event"
// var scheduleApi = baseAddress + "/api/v1/schedule"
// var badfile = "testdata/badfile.xml"
// var goodfile = "testdata/sample_schedule.xml"
// var testschedule dao.ScheduleBlueprint

// func TestMain(m *testing.M) {
// 	go server.Start()
// 	exitVal := m.Run()
// 	go server.Stop()
// 	os.Exit(exitVal)
// }

// func TestGoodParams(test *testing.T) {
// 	//M = M.Init(&c)

// 	goodRequest := "{\"name\": \"kaela\", \"success\": true}"
// 	response, err := http.Post(eventApi, "application/json", strings.NewReader(goodRequest))

// 	assert.Nil(test, err, "Error contacting server")
// 	assert.NotNil(test, response, "Response object is nil")
// 	assert.Equal(test, http.StatusOK, response.StatusCode, "Response http status code not what it should be")

// }
// func TestBadParams(test *testing.T) {

// 	badReqeust := "{}"
// 	response, err := http.Post(eventApi, "application/json", strings.NewReader(badReqeust))

// 	assert.Nil(test, err, "Error contacting server")
// 	assert.NotNil(test, response, "Response object is nil")
// 	assert.Equal(test, http.StatusBadRequest, response.StatusCode, "Response http status code not what it should be")

// }

// func TestGoodSchedule(test *testing.T) {

// 	xfile, err := os.Open(goodfile)
// 	assert.Nil(test, err, "Error opening file")
// 	b, err := ioutil.ReadAll(xfile)

// 	assert.Nil(test, err, "Error getting bytes.")
// 	assert.NotNil(test, xfile, "XML returns nil")
// 	err = xml.Unmarshal(b, &testschedule)
// 	assert.Nil(test, err, "Could not decode bytes.")

// 	response, err := http.NewRequest("PUT", scheduleApi, bytes.NewBuffer(b))
// 	assert.Nil(test, err, "Error getting ready for post")
// 	assert.NotNil(test, response, "Response is nil")
// 	req, err := http.NewRequest("GET", scheduleApi, nil)
// 	assert.Nil(test, err, "Error in GET method")
// 	assert.NotNil(test, req, "Did not receive anything.")
// }

// func TestBadSchedule(test *testing.T) {

// 	xfile, err := os.Open(badfile)
// 	assert.Nil(test, err, "Error opening file")
// 	bytes, err := ioutil.ReadAll(xfile)
// 	assert.Nil(test, err, "Error getting bytes.")
// 	assert.NotNil(test, xfile, "XML returns nil")
// 	err = xml.Unmarshal(bytes, &testschedule)
// 	assert.NotNil(test, err, "Could not decode bytes.")

// }
