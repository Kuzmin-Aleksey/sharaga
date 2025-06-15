package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/failure"
)

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) Save(ctx context.Context, order *entity.Order) error {
	if _, err := r.db.NamedExecContext(ctx, "INSERT INTO orders (creator_id, partner_id, created_at) VALUES (:creator_id, :partner_id, :created_at)", order); err != nil {
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

func (r *OrdersRepo) GetByPartner(ctx context.Context, partnerId int) ([]entity.Order, error) {
	orders := make([]entity.Order, 0)
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM orders WHERE partner_id=?", partnerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orders, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order

		if err := rows.Scan(&order.Id, &order.CreatorId, &order.PartnerId, &order.CreateAt); err != nil {
			return nil, failure.NewInternalError(err.Error())
		}

		orders = append(orders, order)
	}

	return orders, nil
}
