package entity

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Session struct {
	SessionID    string    `json:"session_id"`
	WebsiteID    string    `json:"website_id"`
	Hostname     string    `json:"hostname,omitempty"`
	Browser      string    `json:"browse,omitempty"`
	OS           string    `json:"os,omitempty"`
	Device       string    `json:"device,omitempty"`
	Screen       string    `json:"screen,omitempty"`
	Language     string    `json:"language,omitempty"`
	Country      string    `json:"country,omitempty"`
	Subdivision1 string    `json:"subdivision1,omitempty"`
	Subdivision2 string    `json:"subdivision2,omitempty"`
	City         string    `json:"city,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type SessionClaims struct {
	Session

	OwnerID string `json:"owner_id"`
	VisitID string `json:"visit_id"`
}

type JWTSessionClaims struct {
	jwt.RegisteredClaims
	SessionClaims
}

type UseSessionPayloadData struct {
	Headers map[string][]string
	Body    RequestPayload
	IP      string
	IsLocal bool
}
