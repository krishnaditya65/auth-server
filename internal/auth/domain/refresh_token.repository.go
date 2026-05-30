package domain

import "context"

type RefreshTokenRepository interface {
	Create(
		ctx context.Context,
		token *RefreshToken,
	) error

	GetByHash(
		ctx context.Context,
		hash string,
	) (*RefreshToken, error)

	Revoke(
		ctx context.Context,
		id string,
	) error
}
