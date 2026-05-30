package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims

	TenantID    string   `json:"tid,omitempty"`
	UserID      string   `json:"uid,omitempty"`
	SessionID   string   `json:"sid,omitempty"`
	Email       string   `json:"email,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Nonce       string   `json:"nonce,omitempty"`
}
