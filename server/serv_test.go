package server

import (
	//"net/http"
	//"net/http/httptest"
	"egbitbucket.dtvops.net/deadline/common"
	"log"
	"fmt"
	"net/http"
    "encoding/json"
	"net/http/httptest"
	"testing"
)


var testserv common.DeadlineServer

//-------------------
func testParams(expeced int) (code int) {
  	e := common.Event{}
	testserv.Serv1 = &http.Server{Addr: "8081"}
  if expected == http.StatusOK 
  {
  eventJson, err := json.Marshal(e{Name: "kaela", Success: "false"})
  }else { eventJson, err := json.Marshal(e{Name: "", Success: "false"})  
  
req, err := http.NewRequest("POST", testserv.Serv1.Addr, bytes.NewBuffer(eventJson))

  //post info, wait for server response, 200/400/404
        
	req, err = http.NewRequest("GET", testserv.Serv1.Addr, nil)
	if err != nil {
		fmt.Println("Something went wrong")
		fmt.Println(testserv.Serv1.Addr)
	}
         
	rec := httptest.NewRecorder()
 	eventHandler(rec, req)
	res := rec.Result()
	code = res.StatusCode
	return

}
func TestGoodParams(test *testing.T) {
	//good input, good output
  goodCode := testParams(http.StatusOK)
	if goodCode != http.StatusOK {
		test.Errorf("failed. got %v, expected %v", goodCode, http.StatusOK)
	}
	//assert.Equal(test, 4, 4, "should be equal")

}

func TestBadParams(test *testing.T) {
	//bad input, bad output
	badCode := testParams(httpStatusBadRequest )
	if badCode != http.StatusBadRequest {
		test.Errorf("failed. got %v, expected %v", badCode, http.StatusBadRequest)
	}
}
