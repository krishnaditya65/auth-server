package app

import (
	"context"
	"time"

	authdomain "github.com/krishnaditya65/auth-server/internal/auth/domain"
	authorizationdomain "github.com/krishnaditya65/auth-server/internal/authorization/domain"
	identitydomain "github.com/krishnaditya65/auth-server/internal/identity/domain"

	"github.com/krishnaditya65/auth-server/internal/shared/id"
	"github.com/krishnaditya65/auth-server/internal/shared/password"
	"github.com/krishnaditya65/auth-server/internal/shared/tx"
)

type CreateUserInput struct {
	TenantID    string
	Email       string
	Password    string
	DisplayName string
	RoleName    string
}

type CreateUserUseCase struct {
	txManager tx.Manager

	identityRepo   identitydomain.Repository
	credentialRepo authdomain.Repository

	userRepo identitydomain.UserRepository

	roleRepo     authorizationdomain.RoleRepository
	userRoleRepo authorizationdomain.UserRoleRepository
}

func NewCreateUserUseCase(
	txManager tx.Manager,
	identityRepo identitydomain.Repository,
	credentialRepo authdomain.Repository,
	userRepo identitydomain.UserRepository,
	roleRepo authorizationdomain.RoleRepository,
	userRoleRepo authorizationdomain.UserRoleRepository,
) *CreateUserUseCase {

	return &CreateUserUseCase{
		txManager: txManager,

		identityRepo:   identityRepo,
		credentialRepo: credentialRepo,
		userRepo:       userRepo,

		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
	}
}

func (u *CreateUserUseCase) Execute(
	ctx context.Context,
	input CreateUserInput,
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

			err := u.identityRepo.Create(
				txCtx,
				identity,
			)

			if err != nil {
				return err
			}

			hash, err := password.Hash(
				input.Password,
			)

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

			err = u.credentialRepo.Create(
				txCtx,
				credential,
			)

			if err != nil {
				return err
			}

			user := &identitydomain.User{
				ID:          id.New(),
				TenantID:    input.TenantID,
				IdentityID:  identity.ID,
				Username:    input.Email,
				DisplayName: input.DisplayName,
				Status:      "active",
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			err = u.userRepo.Create(
				txCtx,
				user,
			)

			if err != nil {
				return err
			}

			role, err := u.roleRepo.GetByTenantAndName(
				txCtx,
				input.TenantID,
				input.RoleName,
			)

			if err != nil {
				return err
			}

			err = u.userRoleRepo.AssignRole(
				txCtx,
				&authorizationdomain.UserRole{
					UserID:    user.ID,
					RoleID:    role.ID,
					CreatedAt: now,
				},
			)

			if err != nil {
				return err
			}

			return nil
		},
	)
}
