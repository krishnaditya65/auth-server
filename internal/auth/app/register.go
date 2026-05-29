package app

import (
	"context"
	"time"

	authdomain "github.com/krishnaditya65/auth-server/internal/auth/domain"
	identitydomain "github.com/krishnaditya65/auth-server/internal/identity/domain"

	"github.com/krishnaditya65/auth-server/internal/shared/id"
	"github.com/krishnaditya65/auth-server/internal/shared/password"
	"github.com/krishnaditya65/auth-server/internal/shared/tx"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterUseCase struct {
	txManager      tx.Manager
	identityRepo   identitydomain.Repository
	credentialRepo authdomain.Repository
}

func NewRegisterUseCase(
	txManager tx.Manager,
	identityRepo identitydomain.Repository,
	credentialRepo authdomain.Repository,
) *RegisterUseCase {
	return &RegisterUseCase{
		txManager:      txManager,
		identityRepo:   identityRepo,
		credentialRepo: credentialRepo,
	}
}

func (u *RegisterUseCase) Execute(
	ctx context.Context,
	input RegisterInput,
) error {
	return u.txManager.WithinTransaction(
		ctx,
		func(txCtx context.Context) error {
			now := time.Now().UTC()

			identity := &identitydomain.Identity{
				ID:            id.New(),
				PrimaryEmail:  input.Email,
				EmailVerified: false,
				Status:        "active",
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			if err := u.identityRepo.Create(txCtx, identity); err != nil {
				return err
			}

			hash, err := password.Hash(input.Password)
			if err != nil {
				return err
			}

			credential := &authdomain.Credential{
				ID:             id.New(),
				IdentityID:     identity.ID,
				CredentialType: "password",
				PasswordHash:   hash,
				CreatedAt:      now,
				UpdatedAt:      now,
			}

			return u.credentialRepo.Create(txCtx, credential)
		},
	)
}
