package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/failure"
)

type ProductsRepo struct {
	db *sqlx.DB
}

func NewProductsRepo(db *sqlx.DB) *ProductsRepo {
	return &ProductsRepo{db}
}

func (r *ProductsRepo) Save(ctx context.Context, product *entity.Product) error {
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO products (article, type, name, description, min_price, size_x, size_y, size_z, weight, weight_pack) values (:article, :type, :name, :description, :min_price, :size_x, :size_y, :size_z, :weight, :weight_pack)", product)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	product.Id = int(id)

	return nil
}

func (r *ProductsRepo) FindById(ctx context.Context, id int) (product *entity.Product, err error) {
	product = &entity.Product{}
	if err := r.db.GetContext(ctx, product, "SELECT * FROM products WHERE id=?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.NewNotFoundError(err.Error())
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return product, nil
}

func (r *ProductsRepo) GetAll(ctx context.Context) ([]entity.Product, error) {
	products := make([]entity.Product, 0)
	if err := r.db.SelectContext(ctx, &products, "SELECT * FROM products"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return products, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return products, nil
}

func (r *ProductsRepo) SaveType(ctx context.Context, productType *entity.ProductType) error {
	res, err := r.db.NamedExecContext(ctx, "INSERT INTO product_types (type, k) VALUES (:type, :k)", productType)
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	productTypeId, err := res.LastInsertId()
	if err != nil {
		return failure.NewInternalError(err.Error())
	}
	productType.Id = int(productTypeId)
	return nil
}

func (r *ProductsRepo) UpdateType(ctx context.Context, productType *entity.ProductType) error {
	if _, err := r.db.NamedExecContext(ctx, "UPDATE product_types SET  type=:type, k=:k WHERE id=:id", productType); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (r *ProductsRepo) DeleteType(ctx context.Context, typeId int) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM product_types WHERE id=?", typeId); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (r *ProductsRepo) GetTypes(ctx context.Context) ([]entity.ProductType, error) {
	productTypes := make([]entity.ProductType, 0)
	if err := r.db.SelectContext(ctx, &productTypes, "SELECT * FROM product_types"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return productTypes, nil
		}
		return nil, failure.NewInternalError(err.Error())
	}
	return productTypes, nil
}

func (r *ProductsRepo) Update(ctx context.Context, product *entity.Product) error {
	if _, err := r.db.NamedExecContext(ctx, "UPDATE products SET article=:article, type=:type, name=:name, description=:description, min_price=:min_price, size_x=:size_x, size_y=:size_y, size_z=:size_z, weight=:weight, weight_pack=:weight_pack WHERE id=:id", product); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}

func (r *ProductsRepo) Delete(ctx context.Context, id int) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id); err != nil {
		return failure.NewInternalError(err.Error())
	}
	return nil
}
