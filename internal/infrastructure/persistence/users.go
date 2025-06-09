package persistence

import (
	"github.com/jmoiron/sqlx"
)

type ManagersRepo struct {
	db *sqlx.DB
}

func NewManagersRepo(db *sqlx.DB) *ManagersRepo {
	return &ManagersRepo{db}
}
