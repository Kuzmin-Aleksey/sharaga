package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/authutil"
	"sharaga/pkg/failure"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Save(ctx context.Context, user *entity.User) (err error) {
	user.Password, err = authutil.HashPassword(user.Password)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}

	res, err := r.db.NamedExecContext(ctx, "INSERT INTO users (role, name, email, password) values (:role, :name, :email, :password)", user)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	user.Id = int(id)

	return nil
}

func (r *UsersRepo) FindById(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id=?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return user, nil
}

func (r *UsersRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE email=?", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return user, nil
}

func (r *UsersRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	users := make([]entity.User, 0)
	if err := r.db.SelectContext(ctx, &users, "SELECT * FROM users"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return users, nil
}

func (r *UsersRepo) GetWithPassword(ctx context.Context, username, password string) (*entity.User, error) {
	passwordHash, err := authutil.HashPassword(password)
	if err != nil {
		return nil, failure.NewInternalError(err.Error())
	}

	user := new(entity.User)
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE username=? AND password=?", username, passwordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return user, nil
}

func (r *UsersRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (r *UsersRepo) GetRole(ctx context.Context, userId int) (string, error) {
	var role string
	if err := r.db.QueryRowContext(ctx, "SELECT role FROM users WHERE id=?", userId).Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", failure.NewNotFoundError(err.Error())
		}
		return "", failure.NewInternalError(err.Error())
	}

	return role, nil
}
