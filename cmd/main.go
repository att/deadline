package main

import (
	//"common"
	"deadline/server"
	"fmt"
	//"log"
	//"os/exec"
	//"io"
	"net/http"
)

//read file
//to be moved to server package

//want to configure
func start(d *deadlineServer) {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//m = make(map[string]string)
		var e Event
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
	}
	d.serv1 = &http.Server{Addr: d.serv1.Addr}
	err := d.serv1.ListenAndServe //allow us to start accessing that test server
	if err != nil {
		"We have a problem"
	}
	return
}

func stop(d *deadlineServer) {

	err := d.serv1.Close()
	if err != nil {
	}
	return
}

func main() {
	//read files with things

	dd := deadlineServer{}

	start(&dd)
	fmt.Println("Got past initializing")
	stop(&dd)

}
