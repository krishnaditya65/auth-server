package app

import (
	"context"
	"time"

	authdomain "github.com/krishnaditya65/auth-server/internal/auth/domain"
	identitydomain "github.com/krishnaditya65/auth-server/internal/identity/domain"
	sessiondomain "github.com/krishnaditya65/auth-server/internal/session/domain"

	sharederrors "github.com/krishnaditya65/auth-server/internal/shared/errors"
	"github.com/krishnaditya65/auth-server/internal/shared/id"
	"github.com/krishnaditya65/auth-server/internal/shared/password"
	sharedtoken "github.com/krishnaditya65/auth-server/internal/shared/token"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginUseCase struct {
	identityRepo   identitydomain.Repository
	credentialRepo authdomain.Repository
	sessionRepo    sessiondomain.Repository
}

func NewLoginUseCase(
	identityRepo identitydomain.Repository,
	credentialRepo authdomain.Repository,
	sessionRepo sessiondomain.Repository,
) *LoginUseCase {
	return &LoginUseCase{
		identityRepo:   identityRepo,
		credentialRepo: credentialRepo,
		sessionRepo:    sessionRepo,
	}
}

func (u *LoginUseCase) Execute(
	ctx context.Context,
	input LoginInput,
) (*LoginOutput, error) {
	identity, err := u.identityRepo.GetByEmail(
		ctx,
		input.Email,
	)

	if err != nil {
		return nil, sharederrors.ErrInvalidCredentials
	}

	credential, err := u.credentialRepo.GetByIdentityID(
		ctx,
		identity.ID,
	)

	if err != nil {
		return nil, sharederrors.ErrInvalidCredentials
	}

	valid := password.Verify(
		input.Password,
		credential.PasswordHash,
	)

	if !valid {
		return nil, sharederrors.ErrInvalidCredentials
	}

	refreshToken, err := sharedtoken.GenerateRandom(32)
	if err != nil {
		return nil, err
	}

	refreshHash, err := password.Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	session := &sessiondomain.Session{
		ID:               id.New(),
		IdentityID:       identity.ID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        now.Add(24 * time.Hour),
		CreatedAt:        now,
	}

	err = u.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  "access-token-placeholder",
		RefreshToken: refreshToken,
	}, nil
}
