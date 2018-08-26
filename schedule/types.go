package schedule

import (
	"encoding/xml"
	"time"

	"github.com/att/deadline/common"
)

type Live struct {
	Timing  string         `json:"timing,attr,omitempty" db:"timing"`
	Name    string         `json:"name,attr,omitempty" db:"name"`
	LastRun time.Time      `json:"lastrun"`
	Events  []common.Event `json:"events"`
	Handler common.Handler `json:"handler"`
	Start   common.Node    `json:"-"`
	End     common.Node    `json:"-"`
	Error   common.Node    `json:"-"`
}

type ScheduledHandler struct {
	ScheduleName string `db:"schedulename"`
	Name         string `db:"name"`
	Address      string `db:"address"`
}

type ScheduleManager struct {
	subscriptionTable map[string][]*Live
	ScheduleTable     map[string]*Live
	EvaluationTime    time.Time
}

type innerbytes struct {
	XMLName xml.Name       `xml:"innerbytes"`
	Hander  common.Handler `xml:"handler"`
	Events  []common.Event `xml:"event"`
}
