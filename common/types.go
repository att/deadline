package common
import "encoding/xml"
type Event struct {
	XMLName xml.Name 	  `xml:"event"`
	Name    string            `json:"name" xml:"name,attr"`
	Success bool              `json:"success" xml:"success"`
	Details map[string]string `json:"details,omitempty" xml:"details"`
	ReceiveBy string	  `xml:"receive-by,attr"`
}

type Schedule struct {
	XMLName xml.Name        `xml:"schedule"`
	Handler Handler		`xml:"handler" json:"handler"`
	//XMLName xml.Name 	`xml:"schedule"`
	Timing string		`xml:"timing,attr"`
	Name string 		`xml:"name,attr"`
	Schedule []Event 	`xml:"event"`

}

type Handler struct {

	XMLName xml.Name	`xml:"handler"`
	Name	string		`xml:"name,attr"`
	Address string		`xml:"address"`

}

type Error struct {

	To string 		`xml:"to,attr"`

}
