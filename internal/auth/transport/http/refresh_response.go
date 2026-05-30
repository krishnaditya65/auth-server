package http

type RefreshResponse struct {
	SessionID string `json:"session_id"`

	TenantID string `json:"tenant_id"`

	UserID string `json:"user_id"`

	RefreshToken string `json:"refresh_token"`
}
