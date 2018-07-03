package server 
import (
"net/http"
"net/http/httptest"
"testing"
"fmt"
)
//test for correct input, input at all 
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

    for tc, tp := range cases {
	fmt.Println("Goes through once")
        req, _ := http.NewRequest("GET", "/", nil)
        q := req.URL.Query()
        for k, v := range tp.params {
            q.Add(k, v)
        }
        req.URL.RawQuery = q.Encode()
        rec := httptest.NewRecorder()
        
        res := rec.Result()
        if res.StatusCode != tp.statusCode {
            test.Errorf("`%v` failed, got %v, expected %v", tc, res.StatusCode, tp.statusCode)
        }
    }

}
