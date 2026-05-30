package app

import (
	"context"

	sessiondomain "github.com/krishnaditya65/auth-server/internal/session/domain"
	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
)

type LogoutInput struct {
	SessionID string
}

type LogoutUseCase struct {
	sessionRepo sessiondomain.Repository
}

func NewLogoutUseCase(sessionRepo sessiondomain.Repository) *LogoutUseCase {
	return &LogoutUseCase{sessionRepo: sessionRepo}
}

func (u *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	session, err := u.sessionRepo.GetByID(ctx, input.SessionID)
	if err != nil {
		return sharederrors.ErrNotFound
	}

	if session.RevokedAt != nil {
		return nil
	}

	return u.sessionRepo.Revoke(ctx, session.ID)
}
