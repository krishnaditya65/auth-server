package app

import (
	"context"
	"fmt"

	tenantdomain "github.com/krishnaditya65/auth-server/internal/tenant/domain"
)

type SlugService struct {
	repo tenantdomain.Repository
}

func NewSlugService(
	repo tenantdomain.Repository,
) *SlugService {
	return &SlugService{
		repo: repo,
	}
}

func (s *SlugService) GenerateUniqueSlug(
	ctx context.Context,
	base string,
) (string, error) {

	exists, err := s.repo.ExistsBySlug(
		ctx,
		base,
	)
	if err != nil {
		return "", err
	}

	if !exists {
		return base, nil
	}

	for i := 1; i < 10000; i++ {
		candidate := fmt.Sprintf(
			"%s-%d",
			base,
			i,
		)

		exists, err := s.repo.ExistsBySlug(
			ctx,
			candidate,
		)

		if err != nil {
			return "", err
		}

		if !exists {
			return candidate, nil
		}
	}

	return "", fmt.Errorf(
		"unable to generate unique slug",
	)
}
