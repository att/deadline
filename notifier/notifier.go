package notifier 
import (
	"encoding/json"
	"bytes"
	"net/http"
	"egbitbucket.dtvops.net/deadline/common"
)


func (w Webhook) Send(msg string) {
	str := msg
	jv, err := json.Marshal(str)
	if err != nil {
		common.CheckError(err)
	}
	_ , err = http.Post(w.Addr,"application/json", bytes.NewBuffer(jv))
	common.CheckError(err)


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
	common.Info.Println("Did not give a valid handler.")
	return &Webhook{}
}


