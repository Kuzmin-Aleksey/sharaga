package orders

import (
	"context"
	"fmt"
	"sharaga/internal/domain/aggregate"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/contextx"
)

type Repository interface {
	Save(ctx context.Context, order *aggregate.OrderProducts) error
	Get(ctx context.Context, id int) (*entity.Order, error)
	GetByPartner(ctx context.Context, partnerId int) ([]aggregate.OrderProductInfo, error)
	GetPartnerProductCount(ctx context.Context, partnerId int) (int, error)
	GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error)
}

type ProductProvider interface {
	FindById(ctx context.Context, id int) (product *entity.Product, err error)
}

type Service struct {
	repo     Repository
	products ProductProvider
}

func NewService(repo Repository, products ProductProvider) *Service {
	return &Service{
		repo:     repo,
		products: products,
	}
}

func (s *Service) NewOrder(ctx context.Context, order *aggregate.OrderProducts) error {
	const op = "orders.NewOrder"

	order.Order.CreatorId = int(contextx.GetUserId(ctx))

	order.Order.Price = 0

	for _, orderProd := range order.Products {
		prod, err := s.products.FindById(ctx, orderProd.ProductId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		order.Order.Price += prod.MinPrice * orderProd.Quantity
	}

	discount, err := s.CalcDiscount(ctx, order.Order.PartnerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	order.Order.Price = order.Order.Price - int(float64(order.Order.Price)*float64(discount)/100)

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
	const op = "orders.GetByPartner"

	orders, err := s.repo.GetByPartner(ctx, partnerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return orders, nil
}

func (s *Service) CalcDiscount(ctx context.Context, partnerId int) (int, error) {
	const op = "orders.calcDiscount"

	count, err := s.repo.GetPartnerProductCount(ctx, partnerId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if count > 300000 {
		return 15, nil
	}
	if count > 50000 {
		return 10, nil
	}
	if count > 10000 {
		return 5, nil
	}

	return 0, nil
}
