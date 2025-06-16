package products

import (
	"context"
	"fmt"
	"sharaga/internal/domain/aggregate"
	"sharaga/internal/domain/entity"
)

type Repository interface {
	Save(ctx context.Context, product *entity.Product) error
	GetAllWithType(ctx context.Context) ([]aggregate.ProductWithType, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error

	SaveType(ctx context.Context, productType *entity.ProductType) error
	GetTypes(ctx context.Context) ([]entity.ProductType, error)
	UpdateType(ctx context.Context, productType *entity.ProductType) error
	DeleteType(ctx context.Context, typeId int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) NewProduct(ctx context.Context, product *entity.Product) error {
	const op = "service.NewProduct"
	if err := s.repo.Save(ctx, product); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) GetAllWithType(ctx context.Context) ([]aggregate.ProductWithType, error) {
	const op = "service.GetAllWithType"

	products, err := s.repo.GetAllWithType(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return products, nil
}

func (s *Service) Update(ctx context.Context, product *entity.Product) error {
	const op = "service.Update"
	if err := s.repo.Update(ctx, product); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	const op = "service.Delete"
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) NewType(ctx context.Context, productType *entity.ProductType) error {
	const op = "service.NewType"
	if err := s.repo.SaveType(ctx, productType); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) GetTypes(ctx context.Context) ([]entity.ProductType, error) {
	const op = "service.GetTypes"
	productTypes, err := s.repo.GetTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return productTypes, nil
}

func (s *Service) UpdateType(ctx context.Context, productType *entity.ProductType) error {
	const op = "service.UpdateType"
	if err := s.repo.UpdateType(ctx, productType); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) DeleteType(ctx context.Context, typeId int) error {
	const op = "service.DeleteType"
	if err := s.repo.DeleteType(ctx, typeId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
