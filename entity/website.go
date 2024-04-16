package entity

import "time"

type Website struct {
	WebsiteID string    `json:"website_id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain,omitempty"`
	ShareID   string    `json:"share_id,omitempty"`
	ResetAt   time.Time `json:"reset_at,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	TimeID    string    `json:"time_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type WebsiteEvent struct {
	EventID        string    `json:"event_id"`
	WebsiteID      string    `json:"website_id"`
	SessionID      string    `json:"session_id"`
	VisitID        string    `json:"visit_id"`
	CreatedAt      time.Time `json:"created_at"`
	UrlPath        string    `json:"url_path"`
	UrlQuery       string    `json:"url_query"`
	ReferrerPath   string    `json:"referrer_path"`
	ReferrerQuery  string    `json:"referrer_query"`
	ReferrerDomain string    `json:"referrer_domain"`
	PageTitle      string    `json:"page_title"`
	EventType      int       `json:"event_type"`
	EventName      string    `json:"event_name"`
}
