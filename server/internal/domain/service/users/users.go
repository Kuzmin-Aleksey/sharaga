package users

import (
	"context"
	"fmt"
)

type Repository interface {
	GetRole(ctx context.Context, userId int) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRole(ctx context.Context, userId int) (string, error) {
	const op = "users.GetRole"
	role, err := s.repo.GetRole(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
