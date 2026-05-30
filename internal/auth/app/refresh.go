package app

import (
	"context"
	"time"

	sessiondomain "github.com/krishnaditya65/auth-server/internal/session/domain"
	"github.com/krishnaditya65/auth-server/internal/shared/id"
	sharedtoken "github.com/krishnaditya65/auth-server/internal/shared/token"

	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
)

type RefreshInput struct {
	RefreshToken string
}
type RefreshOutput struct {
	SessionID string

	TenantID string
	UserID   string

	RefreshToken string
}
type RefreshUseCase struct {
	sessionRepo sessiondomain.Repository
}

func NewRefreshUseCase(
	sessionRepo sessiondomain.Repository,
) *RefreshUseCase {
	return &RefreshUseCase{
		sessionRepo: sessionRepo,
	}
}

func (u *RefreshUseCase) Execute(
	ctx context.Context,
	input RefreshInput,
) (*RefreshOutput, error) {

	refreshHash := sharedtoken.Hash(input.RefreshToken)

	session, err := u.sessionRepo.GetByRefreshTokenHash(
		ctx,
		refreshHash,
	)

	if err != nil {
		return nil, sharederrors.ErrUnauthorized
	}

	if session.RevokedAt != nil {
		return nil, sharederrors.ErrUnauthorized
	}

	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, sharederrors.ErrUnauthorized
	}

	err = u.sessionRepo.Revoke(
		ctx,
		session.ID,
	)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := sharedtoken.GenerateRandom(32)

	if err != nil {
		return nil, err
	}

	newRefreshHash := sharedtoken.Hash(newRefreshToken)

	now := time.Now().UTC()

	parentID := session.ID

	newSession := &sessiondomain.Session{
		ID: id.New(),

		TenantID: session.TenantID,

		IdentityID: session.IdentityID,

		UserID: session.UserID,

		RefreshTokenHash: newRefreshHash,

		ParentSessionID: &parentID,

		IPAddress: session.IPAddress,

		UserAgent: session.UserAgent,

		ExpiresAt: now.Add(24 * time.Hour),

		CreatedAt: now,
	}

	err = u.sessionRepo.Create(
		ctx,
		newSession,
	)

	if err != nil {
		return nil, err
	}

	return &RefreshOutput{
		SessionID:    newSession.ID,
		TenantID:     newSession.TenantID,
		UserID:       newSession.UserID,
		RefreshToken: newRefreshToken,
	}, nil
}
