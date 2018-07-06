package server

import (
	"net/http"
  "egbitbucket.dtvops.net/deadline/common"
  "log"
  "io/ioutil"
  
)

func eventHander( w http.ResponseWriter, r *http.Request)
{
  e := common.Event{} 
  jsn, err := ioutil.ReadAll(r.Body)
  if err != nil { log.Fatal("Error reading the body", err }
  err = json.Unmarshal(jsn, &e)
  if err != nil { log.Fatal("Could not decode" ,err) }
  //see if we successfully got the struct
  log.Printf("Received the following information: %v\n", e)   
  checkParams(e)
  //return 200/400/404
  
}

func checkParams(event common.Event) () {
  
  //type, strings, xml, etc.
   
}
                            //func getQueryString ()(event common.Event) {
                            //}
