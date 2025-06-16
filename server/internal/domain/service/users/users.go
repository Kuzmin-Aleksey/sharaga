package users

import (
	"context"
	"fmt"
	"sharaga/internal/domain/entity"
)

type Repository interface {
	Save(ctx context.Context, user *entity.User) (err error)
	Update(ctx context.Context, user *entity.User) (err error)
	GetAll(ctx context.Context) ([]entity.User, error)
	GetRole(ctx context.Context, userId int) (string, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewUser(ctx context.Context, user *entity.User) error {
	const op = "users.NewUser"
	if err := s.repo.Save(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) UpdateUser(ctx context.Context, user *entity.User) error {
	const op = "users.UpdateUser"
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userId int) error {
	const op = "users.DeleteUser"
	if err := s.repo.Delete(ctx, userId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]entity.User, error) {
	const op = "users.GetAll"
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return users, nil
}

func (s *Service) GetRole(ctx context.Context, userId int) (string, error) {
	const op = "users.GetRole"
	role, err := s.repo.GetRole(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
