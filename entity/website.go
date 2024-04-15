package entity

import "time"

type Website struct {
	ID        string    `json:"website_id"`
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
