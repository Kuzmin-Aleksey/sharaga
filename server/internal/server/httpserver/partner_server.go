package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/failure"
	"sharaga/pkg/rest"
	"strconv"
)

type partnerService interface {
	NewPartner(ctx context.Context, partner *entity.Partner) error
	GetAll(ctx context.Context) ([]entity.Partner, error)
	Update(ctx context.Context, partner *entity.Partner) error
	Delete(ctx context.Context, partnerId int) error
}

type PartnerServer struct {
	partnerService partnerService
}

func NewPartnerServer(partnerService partnerService) *PartnerServer {
	return &PartnerServer{
		partnerService: partnerService,
	}
}

func (p *PartnerServer) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partner := new(entity.Partner)

	if err := json.NewDecoder(r.Body).Decode(partner); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := p.partnerService.NewPartner(ctx, partner); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.IdResponse{
		Id: partner.Id,
	}, http.StatusOK)
}

func (p *PartnerServer) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partners, err := p.partnerService.GetAll(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, partners, http.StatusOK)
}

func (p *PartnerServer) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	partner := new(entity.Partner)
	if err := json.NewDecoder(r.Body).Decode(partner); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := p.partnerService.Update(ctx, partner); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *PartnerServer) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	partnerId, err := strconv.Atoi(r.FormValue("partner_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError("invalid partner id: "+r.FormValue("partner_id")))
		return
	}

	if err := p.partnerService.Delete(ctx, partnerId); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
