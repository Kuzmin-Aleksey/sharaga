package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sharaga/internal/domain/aggregate"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/contextx"
	"sharaga/pkg/failure"
	"sharaga/pkg/logx"
)

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) Save(ctx context.Context, order *aggregate.OrderProducts) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return failure.NewInternalError(err.Error())
	}

	if _, err := tx.NamedExecContext(ctx, "INSERT INTO orders (creator_id, partner_id, created_at) VALUES (:creator_id, :partner_id, :created_at)", order.Order); err != nil {
		return failure.NewInternalError(err.Error())
	}

	for _, product := range order.Products {
		if _, err := tx.NamedExecContext(ctx, "INSERT INTO order_product (order_id, product_id, quantity) VALUES (:order_id, :product_id, :quantity)", product); err != nil {
			if err := tx.Rollback(); err != nil {
				contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "save order: rollback failed: ", logx.Error(err))
			}
			return failure.NewInternalError(err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return failure.NewInternalError(err.Error())
	}

	return nil
}

func (r *OrdersRepo) Get(ctx context.Context, id int) (*entity.Order, error) {
	order := new(entity.Order)
	if err := r.db.GetContext(ctx, order, "SELECT * FROM orders WHERE id=?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return order, nil
}

func (r *OrdersRepo) GetByPartner(ctx context.Context, partnerId int) ([]aggregate.OrderProductInfo, error) {
	ordersWithProducts := make([]aggregate.OrderProductInfo, 0)
	orders := make([]entity.Order, 0)
	if err := r.db.SelectContext(ctx, &orders, "SELECT * FROM orders WHERE partner_id=?", partnerId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ordersWithProducts, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}

	for _, order := range orders {
		productInfo := make([]aggregate.ProductQuantity, 0)

		if err := r.db.SelectContext(ctx, productInfo, `
			SELECT p.id AS id, p.article AS article, p.type AS type, p.name AS name, p.description AS description, p.min_price AS min_price, p.size_x AS size_x, p.size_y AS size_y, p.size_z AS size_z, p.weight AS weight, p.weight_pack AS weight_pack, op.quantity AS quantity 
 			FROM order_product op
			INNER JOIN products p ON op.product_id = p.id
`, order); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, failure.NewInternalError(err.Error())
			}
		}

		ordersWithProducts = append(ordersWithProducts, aggregate.OrderProductInfo{
			Order:    order,
			Products: productInfo,
		})
	}

	return ordersWithProducts, nil
}

func (r *OrdersRepo) GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error) {
	ordersWithProducts := make([]aggregate.OrderProductInfo, 0)
	orders := make([]entity.Order, 0)
	if err := r.db.SelectContext(ctx, &orders, "SELECT * FROM orders"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ordersWithProducts, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}

	for _, order := range orders {
		productInfo := make([]aggregate.ProductQuantity, 0)

		if err := r.db.SelectContext(ctx, productInfo, `
			SELECT p.id AS id, p.article AS article, p.type AS type, p.name AS name, p.description AS description, p.min_price AS min_price, p.size_x AS size_x, p.size_y AS size_y, p.size_z AS size_z, p.weight AS weight, p.weight_pack AS weight_pack, op.quantity AS quantity 
 			FROM order_product op
			INNER JOIN products p ON op.product_id = p.id
`, order); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, failure.NewInternalError(err.Error())
			}
		}

		ordersWithProducts = append(ordersWithProducts, aggregate.OrderProductInfo{
			Order:    order,
			Products: productInfo,
		})
	}

	return ordersWithProducts, nil
}
