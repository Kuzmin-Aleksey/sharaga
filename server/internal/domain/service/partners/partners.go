package partners

import (
	"context"
	"fmt"
	"sharaga/internal/domain/entity"
)

type Repository interface {
	Save(ctx context.Context, partner *entity.Partner) error
	GetAll(ctx context.Context) ([]entity.Partner, error)
	Update(ctx context.Context, partner *entity.Partner) error
	Delete(ctx context.Context, partnerId int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewPartner(ctx context.Context, partner *entity.Partner) error {
	const op = "partners.NewPartner"
	if err := s.repo.Save(ctx, partner); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]entity.Partner, error) {
	const op = "partners.GetAll"
	partners, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return partners, nil
}

func (s *Service) Update(ctx context.Context, partner *entity.Partner) error {
	const op = "partners.Update"
	if err := s.repo.Update(ctx, partner); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, partnerId int) error {
	const op = "partners.Delete"
	if err := s.repo.Delete(ctx, partnerId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
