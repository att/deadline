package main

import "io"
import "bytes"
//import "net"
import "fmt"
//import "bufio"
import "os"
//import "log"
//import "time"
import "encoding/json"
import "net/http"
func randomize() {

//will randomize values for potential test cases in the future 




}


func main() {

  // connect to this socket
//  conn, err:= net.Dial("tcp", "localhost:8081")
//	if err != nil { log.Fatal(err)}
//	if conn != nil {}
fmt.Println("Starting client..")

//make a struct
type event struct {

	Name string
	Success bool
	Details map[string]string 




	}	

m := make(map[string]string) //makes an empty map 


e := event{"Kaela", false, m}
buf := new(bytes.Buffer) 
 // for { 
   //try gob?
json.NewEncoder(buf).Encode(e)
// fmt.Println("Before post")
resp, _ := http.Post(":8081","application/json;charset=utf-8",buf) 


_, err := io.Copy(os.Stdout, resp.Body) 
if  err == io.EOF { fmt.Println("Well it copied")
}

  
    
    
//fmt.Fprintf(conn, e.Name + "\n")
  
//reply
  //  message, _ := bufio.NewReader(conn).ReadString('\n')
 //  fmt.Print("Message from server: "+message) 
//	time.Sleep(500*time.Millisecond)
 //}
}
