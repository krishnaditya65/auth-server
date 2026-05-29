package principal

type Principal struct {
	SessionID string

	TenantID string
	UserID   string

	IdentityID string

	Roles []string
}
