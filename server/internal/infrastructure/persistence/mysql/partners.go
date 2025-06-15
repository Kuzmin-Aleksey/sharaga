package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/failure"
)

type PartnersRepo struct {
	db *sqlx.DB
}

func NewPartnersRepo(db *sqlx.DB) *PartnersRepo {
	return &PartnersRepo{db: db}
}

func (r *PartnersRepo) Save(ctx context.Context, partner *entity.Partner) error {
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO partners (type, name, director, email, phone, address, inn, rating) Values(:type, :name, :director, :email, :phone, :address, :inn, :rating)", partner)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	partner.Id = int(id)
	return nil
}

func (r *PartnersRepo) FindById(ctx context.Context, id int) (*entity.Partner, error) {
	partner := &entity.Partner{}
	if err := r.db.GetContext(ctx, partner, "SELECT * FROM partners WHERE id=?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return partner, nil
}

func (r *PartnersRepo) GetAll(ctx context.Context) ([]entity.Partner, error) {
	partners := make([]entity.Partner, 0)
	if err := r.db.SelectContext(ctx, &partners, "SELECT * FROM partners"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return partners, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return partners, nil
}

func (r *PartnersRepo) Update(ctx context.Context, partner *entity.Partner) error {
	if _, err := r.db.NamedExecContext(ctx, "UPDATE partners SET type=:type, name=:name, director=:director, email=:email, phone=:phone, address=:address, inn=:inn, rating=:rating WHERE id=:id", partner); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (r *PartnersRepo) Delete(ctx context.Context, id int) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM partners WHERE id=?", id); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}
