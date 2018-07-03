package main
import (
//"common"
//"server"
"fmt"
//"log"
//"os/exec"
//"io"
"net/http"
)

//read file 
//to be moved to server package 
type deadlineServer struct {

serv1 *http.Server
address string 
//handler?
//handler *http.Handler 
} 


//want to configure 
func start (d *deadlineServer) {

d.serv1 = &http.Server{Addr: d.address}
err := d.serv1.ListenAndServe //allow us to start accessing that test server 
if err != nil {} 
return 
}

func stop ( d *deadlineServer) {

	err := d.serv1.Close()
	if err != nil {} 
	return
}

func main(){
//read files with things 

dd := deadlineServer {
address: ":8081",
}

start(&dd)
fmt.Println("Got past initializing")
stop(&dd)



}

