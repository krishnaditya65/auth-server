package middleware

import (
	"net/http"
	"time"

	authdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"
	identitydomain "github.com/krishnaditya65/auth-server/internal/identity/domain"
	sessiondomain "github.com/krishnaditya65/auth-server/internal/session/domain"
	authctx "github.com/krishnaditya65/auth-server/internal/shared/context"
	"github.com/krishnaditya65/auth-server/internal/shared/principal"
)

type AuthenticationMiddleware struct {
	sessionRepo  sessiondomain.Repository
	identityRepo identitydomain.Repository
	userRoleRepo authdomain.UserRoleRepository
}

func NewAuthenticationMiddleware(
	sessionRepo sessiondomain.Repository,
	identityRepo identitydomain.Repository,
	userRoleRepo authdomain.UserRoleRepository,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		sessionRepo:  sessionRepo,
		identityRepo: identityRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (m *AuthenticationMiddleware) Authenticate(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			sessionID := r.Header.Get(
				"X-Session-ID",
			)

			if sessionID == "" {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			session, err := m.sessionRepo.GetByID(
				r.Context(),
				sessionID,
			)

			if err != nil {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			if session.RevokedAt != nil {
				http.Error(
					w,
					"session revoked",
					http.StatusUnauthorized,
				)
				return
			}

			if time.Now().UTC().After(
				session.ExpiresAt,
			) {
				http.Error(
					w,
					"session expired",
					http.StatusUnauthorized,
				)
				return
			}

			identity, err := m.identityRepo.GetByID(
				r.Context(),
				session.IdentityID,
			)

			if err != nil {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			roles, err := m.userRoleRepo.GetRolesForUser(
				r.Context(),
				session.UserID,
			)

			if err != nil {
				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)
				return
			}

			p := &principal.Principal{
				SessionID: session.ID,

				IdentityID: session.IdentityID,

				TenantID: session.TenantID,
				UserID:   session.UserID,

				Email: identity.PrimaryEmail,

				Roles: roles,
			}

			ctx := authctx.WithPrincipal(
				r.Context(),
				p,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		},
	)
}
