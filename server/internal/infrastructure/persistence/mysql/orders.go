package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
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

	res, err := tx.NamedExecContext(ctx, "INSERT INTO orders (creator_id, partner_id, create_at, price) VALUES (:creator_id, :partner_id, :create_at, :price)", order.Order)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}

	orderId, err := res.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "save order: rollback failed: ", logx.Error(err))
		}
		return failure.NewInternalError(err.Error())
	}
	order.Order.Id = int(orderId)

	for i, product := range order.Products {
		order.Products[i].OrderId = order.Order.Id
		product.OrderId = order.Order.Id
		log.Println("insert product:", product)

		if _, err := tx.NamedExecContext(ctx, "INSERT INTO order_product (order_id, product_id, quantity, price) VALUES (:order_id, :product_id, :quantity, :price)", product); err != nil {
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
			SELECT p.id AS id, p.article AS article, p.type AS type, p.name AS name, p.description AS description, p.min_price AS min_price, p.size_x AS size_x, p.size_y AS size_y, p.size_z AS size_z, p.weight AS weight, p.weight_pack AS weight_pack, op.quantity, op.price AS quantity 
 			FROM order_product op
			INNER JOIN products p ON p.id = op.product_id
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

func (r *OrdersRepo) GetPartnerProductCount(ctx context.Context, partnerId int) (int, error) {
	var count int

	if err := r.db.GetContext(ctx, &count, `
			SELECT SUM(op.quantity) AS count
 			FROM order_product op
			INNER JOIN orders o ON op.order_id = o.id
			WHERE o.partner_id=?`, partnerId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
	}

	return count, nil
}

func (r *OrdersRepo) GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error) {
	ordersWithProducts := make([]aggregate.OrderProductInfo, 0)
	orders := make([]entity.Order, 0)
	if err := r.db.SelectContext(ctx, &orders, "SELECT * FROM orders ORDER BY create_at DESC"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ordersWithProducts, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}

	for _, order := range orders {
		productInfo := make([]aggregate.ProductQuantity, 0)

		if err := r.db.SelectContext(ctx, &productInfo, `
			SELECT p.id AS id, p.article AS article, p.type AS type, p.name AS name, p.description AS description, p.min_price AS min_price, p.size_x AS size_x, p.size_y AS size_y, p.size_z AS size_z, p.weight AS weight, p.weight_pack AS weight_pack, op.quantity AS quantity , op.price as price
 			FROM order_product op
			INNER JOIN products p ON op.product_id = p.id
			WHERE op.order_id = ? 
`, order.Id); err != nil {
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
