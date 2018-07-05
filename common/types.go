package common
import "net/http"

type Event struct {

	Name string
        Success bool
        Details map[string]string

}


type DeadlineServer struct {
        Serv1 *http.Server
}

