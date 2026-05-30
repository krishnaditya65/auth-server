package domain

import (
	"context"
	"time"
)

type AuthorizationCode struct {
	Code                string    `json:"code"`
	ClientID            string    `json:"client_id"`
	TenantID            string    `json:"tenant_id"`
	IdentityID          string    `json:"identity_id"`
	UserID              string    `json:"user_id"`
	Scopes              []string  `json:"scopes"`
	RedirectURI         string    `json:"redirect_uri"`
	CodeChallenge       string    `json:"code_challenge"`
	CodeChallengeMethod string    `json:"code_challenge_method"`
	Nonce               string    `json:"nonce"`
	ExpiresAt           time.Time `json:"expires_at"`
}

type CodeStore interface {
	Store(ctx context.Context, code *AuthorizationCode) error
	Consume(ctx context.Context, code string) (*AuthorizationCode, error)
}
