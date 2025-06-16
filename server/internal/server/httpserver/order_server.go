package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sharaga/internal/domain/aggregate"
	"sharaga/pkg/failure"
	"sharaga/pkg/rest"
	"strconv"
)

type orderService interface {
	NewOrder(ctx context.Context, order *aggregate.OrderProducts) error
	GetAll(ctx context.Context) ([]aggregate.OrderProductInfo, error)
	GetByPartner(ctx context.Context, partnerId int) ([]aggregate.OrderProductInfo, error)
	CalcDiscount(ctx context.Context, partnerId int) (int, error)
}

type OrderServer struct {
	orderService orderService
}

func NewOrderServer(orderService orderService) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}

func (s *OrderServer) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	order := new(aggregate.OrderProducts)

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.orderService.NewOrder(ctx, order); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.IdResponse{
		Id: order.Order.Id,
	}, http.StatusOK)
}

func (s *OrderServer) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := s.orderService.GetAll(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, orders, http.StatusOK)
}

func (s *OrderServer) GetByPartner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partnerId, err := strconv.Atoi(r.FormValue("partner_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError("invalid partner id: "+r.FormValue("partner_id")))
		return
	}

	orders, err := s.orderService.GetByPartner(ctx, partnerId)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, orders, http.StatusOK)
}

func (s *OrderServer) CalcDiscount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partnerId, err := strconv.Atoi(r.FormValue("partner_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError("invalid partner id: "+r.FormValue("partner_id")))
		return
	}

	discount, err := s.orderService.CalcDiscount(ctx, partnerId)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.DiscountResponse{
		Discount: discount,
	}, http.StatusOK)
}
