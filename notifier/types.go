package notifier 
import (
	"net/http"
)
//basically making a set 


type NotifyHandler interface {
	Send(string) 
}

type TypeHandler struct {
	Name string
}
type Webhook struct {
	TH 	 TypeHandler
	Addr string
	Handler http.Handler

}
