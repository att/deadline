package server 
import (
"net/http"
"net/http/httptest"
"testing"
)
func TestParams (test *testing.T) {

//come up with test cases


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
        req, _ := http.NewRequest("GET", "/", nil)
        q := req.URL.Query()
        for k, v := range tp.params {
            q.Add(k, v)
        }
        req.URL.RawQuery = q.Encode()
        rec := httptest.NewRecorder()
        //root(rec, req)
        res := rec.Result()
        if res.StatusCode != tp.statusCode {
            test.Errorf("`%v` failed, got %v, expected %v", tc, res.StatusCode, tp.statusCode)
        }
    }






}



//parameters 
