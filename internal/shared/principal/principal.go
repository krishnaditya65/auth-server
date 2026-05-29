package principal

type Principal struct {
	SessionID string

	IdentityID string

	TenantID string
	UserID   string

	Email string

	Roles []string

	Permissions []string
}
