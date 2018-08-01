package notifier 
import (
	"encoding/json"
	"bytes"
	"log"
	"net/http"
)
//way of keeping track if there is a handler made for it

func (w Webhook) Send(msg string) {
	str := msg
	jv, err := json.Marshal(str)
	if err != nil {
		log.Println("Could not encode.")
	}
	_ , err = http.Post(w.Addr+"/api/v1/msg","application/json", bytes.NewBuffer(jv))

	if err != nil {
		log.Println(err)
	}


}
//make a different function for each type 
func NewNotifyHandler(nh string) NotifyHandler{

	switch nh {
	case "WEBHOOK":	
		//add to map 
		//s := handlers[nh]
		w := &Webhook{
			Addr: "http://localhost:8081",
		}
		//s = append(s, w)
		//handlers[nh] = s
		return w
	}
	log.Println("Did not give a valid handler.")
	return &Webhook{}
}


