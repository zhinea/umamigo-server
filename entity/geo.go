package entity

type GeoLocation struct {
	Country      string `json:"country,omitempty"`
	Subdivision1 string `json:"subdivision1,omitempty"`
	Subdivision2 string `json:"subdivision2,omitempty"`
	City         string `json:"city,omitempty"`
}

type GeoAgent struct {
	Browser string `json:"browser,omitempty"`
	OS      string `json:"os,omitempty"`
	Device  string `json:"device,omitempty"`
}

type GeoClientInfo struct {
	UserAgent    string `json:"user_agent,omitempty"`
	Browser      string `json:"browser,omitempty"`
	OS           string `json:"os,omitempty"`
	IP           string `json:"ip,omitempty"`
	Country      string `json:"country,omitempty"`
	Subdivision1 string `json:"subdivision1,omitempty"`
	Subdivision2 string `json:"subdivision2,omitempty"`
	City         string `json:"city,omitempty"`
	Device       string `json:"device,omitempty"`
}
