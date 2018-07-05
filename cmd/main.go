package main

import (
	//"common"
	"encoding/json"
	"deadline/common"
	"fmt"
	//"log"
	//"os/exec"
	//"io"
	"net/http"
)

//read file
//to be moved to server package

//want to configure
func start(d *common.DeadlineServer) {


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//m = make(map[string]string)
		var e common.Event
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
	d.Serv1 = &http.Server{Addr: d.Serv1.Addr}
	err := d.Serv1.ListenAndServe //allow us to start accessing that test server
	if err != nil {
		fmt.Println("We have a problem")
	}
	return
}

func stop(d *common.DeadlineServer) {

	err := d.Serv1.Close()
	if err != nil {
	}
	return
}

func main() {
	//read files with things

	dd := common.DeadlineServer{}

	start(&dd)
	fmt.Println("Got past initializing")
	stop(&dd)

}
