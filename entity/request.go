package entity

type RequestPayload struct {
	Events []struct {
		Type    string `json:"type"`
		Payload struct {
			Data     any    `json:"data,omitempty"`
			Referrer string `json:"referrer,omitempty"`
			Title    string `json:"title"`
			Url      string `json:"url"`
			Name     string `json:"name,omitempty"`
			Tag      string `json:"tag,omitempty"`
		} `json:"payload"`
	} `json:"events"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip,omitempty"`
	Screen   string `json:"screen"`
	ID       string `json:"id"`
	Language string `json:"language,omitempty"`
}
