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
func NewNotifyHandler(handlerType string,addr string) NotifyHandler{

	switch handlerType {
	case "WEBHOOK":	

		w := &Webhook{
			Addr: addr,
		}
		w.TH.Name = handlerType

		return w
	}
	common.Info.Println("Did not give a valid handler.")
	return &Webhook{}
}


