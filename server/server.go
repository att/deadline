package server

import (
  "net/http"
  "egbitbucket.dtvops.net/deadline/common"
  "log"
  "io/ioutil"
  "encoding/json"
  
)

func eventHander( w http.ResponseWriter, r *http.Request)
{
 // if r.URL.Path != "8081" -- wrong
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
  s := event.Name.(string)
  b := event.Success.(bool)
  //panics if not proper type 
  //check string length etc. 
}
                            //func getQueryString ()(event common.Event) {
                            //}
