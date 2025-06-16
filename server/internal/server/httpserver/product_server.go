package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sharaga/internal/domain/aggregate"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/failure"
	"sharaga/pkg/rest"
	"strconv"
)

type productService interface {
	NewProduct(ctx context.Context, product *entity.Product) error
	GetAllWithType(ctx context.Context) ([]aggregate.ProductWithType, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error
	NewType(ctx context.Context, productType *entity.ProductType) error
	GetTypes(ctx context.Context) ([]entity.ProductType, error)
	UpdateType(ctx context.Context, productType *entity.ProductType) error
	DeleteType(ctx context.Context, typeId int) error
}

type ProductServer struct {
	productService productService
}

func NewProductServer(productService productService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	product := new(entity.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.productService.NewProduct(ctx, product); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.IdResponse{
		Id: product.Id,
	}, http.StatusOK)
}

func (s *ProductServer) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := s.productService.GetAllWithType(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, products, http.StatusOK)
}

func (s *ProductServer) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	product := new(entity.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.productService.Update(ctx, product); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ProductServer) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productId, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError("invalid product id: "+r.FormValue("product_id")))
		return
	}

	if err := s.productService.Delete(ctx, productId); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ProductServer) NewType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productType := new(entity.ProductType)
	if err := json.NewDecoder(r.Body).Decode(productType); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.productService.NewType(ctx, productType); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.IdResponse{
		Id: productType.Id,
	}, http.StatusOK)
}

func (s *ProductServer) GetTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productTypes, err := s.productService.GetTypes(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, productTypes, http.StatusOK)
}

func (s *ProductServer) UpdateType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productType := new(entity.ProductType)
	if err := json.NewDecoder(r.Body).Decode(productType); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.productService.UpdateType(ctx, productType); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ProductServer) DeleteType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productTypeId, err := strconv.Atoi(r.FormValue("product_type_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.productService.DeleteType(ctx, productTypeId); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
