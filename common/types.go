package common

type Event struct {
	Name    string            `json:"name"`
	Success bool              `json:"success"`
	Details map[string]string `json:"details,omitempty"`
}
