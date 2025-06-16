package orders

import (
	"context"
	"fmt"
	"sharaga/internal/domain/aggregate"
	"sharaga/internal/domain/entity"
)

type Repository interface {
	Save(ctx context.Context, order *aggregate.OrderProducts) error
	Get(ctx context.Context, id int) (*entity.Order, error)
	GetByPartner(ctx context.Context, partnerId int) ([]aggregate.OrderProductInfo, error)
	GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewOrder(ctx context.Context, order *aggregate.OrderProducts) error {
	const op = "orders.NewOrder"

	if err := s.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error) {
	const op = "orders.GetAll"

	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return orders, nil
}

func (s *Service) GetByPartner(ctx context.Context, partnerId int) ([]aggregate.OrderProductInfo, error) {
	const op = "orders.GetAll"

	orders, err := s.repo.GetByPartner(ctx, partnerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}
