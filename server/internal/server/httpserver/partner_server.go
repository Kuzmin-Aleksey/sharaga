package httpserver

type partnerService interface {
}

type PartnerServer struct {
	partnerService partnerService
}

func NewPartnerServer(partnerService partnerService) *PartnerServer {
	return &PartnerServer{
		partnerService: partnerService,
	}
}
