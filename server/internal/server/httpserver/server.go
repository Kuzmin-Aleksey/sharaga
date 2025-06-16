package httpserver

type Server struct {
	userServer    *UserServer
	productServer *ProductServer
	partnerServer *PartnerServer
	orderServer   *OrderServer
	authServer    *AuthServer

	authMw *MwAuth
}

func NewServer(
	userServer *UserServer,
	productServer *ProductServer,
	partnerServer *PartnerServer,
	orderServer *OrderServer,
	authServer *AuthServer,
	authMw *MwAuth,

) *Server {
	var h = &Server{
		userServer:    userServer,
		productServer: productServer,
		partnerServer: partnerServer,
		orderServer:   orderServer,
		authServer:    authServer,

		authMw: authMw,
	}

	return h
}
