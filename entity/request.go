package entity

type RequestPayload struct {
	Events   []EventRequest `json:"events"`
	Hostname string         `json:"hostname"`
	IP       string         `json:"ip,omitempty"`
	Screen   string         `json:"screen"`
	ID       string         `json:"id"`
	Language string         `json:"language,omitempty"`
}
