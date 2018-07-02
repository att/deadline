package server

import "net/http"
import "encoding/json"
//import "net"
import "fmt"
import "log"
//import "bufio"
//import "strings" // only needed below for sample processing

func main() {

type event struct {

        Name string
        Success bool
        Details map[string]string
        }



  fmt.Println("Launching server...")

  // listen on all interfaces
//  ln, err := net.Listen("tcp", ":8081")
//	if err != nil {
//fmt.Println("We might have a problem")
//log.Fatal(err)
//}
//fmt.Println("We get here")
  // accept connection on port
//  conn, err := ln.Accept()
//if (err != nil) {
//log.Fatal(err)

//}
//fmt.Print("We getto after the accept ")
  // run loop forever (or until ctrl-c)
//if conn != nil {} 
// for {
    // will listen for message to process ending in newline (\n)
    //message, _ := bufio.NewReader(conn).ReadString('\n')
    //output message received

   // fmt.Print("Message Received:", string(message))

    // sample process for string received
    //newmessage := strings.ToUpper(message)
    // send new string back to client
    //send struct back to client 
    //conn.Write([]byte(newmessage + "\n"))

//  }
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var e event 
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Println(e.Name)
	})
	log.Fatal(http.ListenAndServe(":8081", nil))


}
