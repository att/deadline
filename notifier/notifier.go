package notifier 
import (
	"encoding/json"
	"bytes"
	"log"
	"net/http"
)


func (w Webhook) Send(msg string) {
	str := msg
	jv, err := json.Marshal(str)
	if err != nil {
		log.Println("Could not encode.")
	}
	_ , err = http.Post(w.Addr,"application/json", bytes.NewBuffer(jv))

	if err != nil {
		log.Println(err)
	}


}
//make a different function for each type 
func NewNotifyHandler(nh string,addr string) NotifyHandler{

	switch nh {
	case "WEBHOOK":	

		w := &Webhook{
			Addr: addr,
		}
		w.TH.Name = nh

		return w
	}
	log.Println("Did not give a valid handler.")
	return &Webhook{}
}


