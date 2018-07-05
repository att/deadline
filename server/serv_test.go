package server

import (
	//"net/http"
	//"net/http/httptest"
	"egbitbucket.dtvops.net/deadline/common"
	//"log"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

//test for correct input, input at all
/*

func TestParams (test *testing.T) {



 cases := map[string]struct {
        params map[string]string
        statusCode int
    }{
        "good params": {
            map[string]string{
                "first": "jeff-test", "second": "true",
            },
            http.StatusOK,
        },
        "without params": {
            map[string]string{},
            http.StatusBadRequest,
        },
    }



}

*/
var testserv common.DeadlineServer

//-------------------
func testParams() (code int) {
	testserv.Serv1 = &http.Server{Addr: "8081"}
	fmt.Println("Get here")
	//for tc, tp := range cases {
	//http.HandleFunc("/", nil)
	//log.Fatal(http.ListenAndServe(":8081", nil))
	req, err := http.NewRequest("GET", testserv.Serv1.Addr, nil)
	if err != nil {
		fmt.Println("Something went wrong")
		fmt.Println(testserv.Serv1.Addr)
	}
	q := req.URL.Query()
	//	for k, v := range tp.params {
	q.Add("12345", "67890")
	//	} /adds string and key
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	//recooooords response
	res := rec.Result()
	code = res.StatusCode
	return

	//	}
	//could have made test tables to avoid duplicate code
}
func TestGoodParams(test *testing.T) {
	//good input, good output
	goodCode := testParams()
	if goodCode != http.StatusOK {
		test.Errorf("failed. got %v, expected %v", goodCode, http.StatusOK)
	}
	//url := serv.URL + ""/deadline/" + something
	//want to take the URL from ttttest server

	//assert.Equal(test, 4, 4, "should be equal")

}

func TestBadParams(test *testing.T) {
	//bad input, bad output
	badCode := testParams()
	if badCode != http.StatusBadRequest {
		test.Errorf("failed. got %v, expected %v", badCode, http.StatusBadRequest)
	}
}
