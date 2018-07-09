package common
import "encoding/xml"
type Event struct {
	XMLName xml.Name `xml:"event"`
	Name    string            `json:"name" xml:"name"`
	Success bool              `json:"success" xml:"success"`
	Details map[string]string `json:"details,omitempty" xml:"details"`
}

type Schedule struct {

	XMLName xml.Name `xml:"users"`
	Schedule []Event `xml:"event"`

}
