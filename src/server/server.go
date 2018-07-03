package server

import  (
"net/http"
"encoding/json"

"fmt"
"log"

)
func main() {
  fmt.Println("Launching server...")


type Event struct {

        Name string	`json:"name"` 
        Success bool	`json:"success"` 	
        Details map[string]string   `json:"details,omitempty"` 
        }




//fmt.Println("Does it even go in here")

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
	})
	log.Fatal(http.ListenAndServe(":8081", nil))


}
