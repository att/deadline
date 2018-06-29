//server that can take in 
/*
name: String
success:bool
details: map[string] string 
*/


package main

import (
"fmt"
"net"
"log"
"time"
"encoding/json"



)
getInf() {



}
func main () {

type commData struct {

ID string
success bool 
details map[string]string]


}

type ( 
	Item struct {
	Name string			'json: name, omitempty'
	Success bool			'json: bool, omitempty'
	Details map[string]string	'json: map, omitempty'

	}


//resp, err := http.Get("http://google.com")
//if resp != nil {}
//if err != nil  {} 
//crate dummy server that spits out information every 5 seconds 

ln,err := net.Listen("tcp, ":602")
if err != nil {} 

//until the connection is no longer valid 

for { 
conn, err := ln.Accept()
if err != nil {}
//go handleConnection(conn)
}

//we will print the struct 
fmt.Println("Compiling. \n")
log.Fatal(http.ListenAndServe(:8080, nil)


}
