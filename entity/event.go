package entity

type EventRequest struct {
	Type    string `json:"type"`
	Payload struct {
		Data     any    `json:"data,omitempty"`
		Referrer string `json:"referrer,omitempty"`
		Title    string `json:"title"`
		Url      string `json:"url"`
		Name     string `json:"name,omitempty"`
		Tag      string `json:"tag,omitempty"`
		T        int    `json:"t,omitempty"`
	} `json:"payload"`
}

type EventCreationPayload struct {
	SessionClaims
	WebsiteEvent
	EventData any `json:"event_data,omitempty"`
	T         int
}
