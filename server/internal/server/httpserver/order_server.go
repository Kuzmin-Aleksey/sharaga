package httpserver

type orderService interface {
}

type OrderServer struct {
	orderService orderService
}

func NewOrderServer(orderService orderService) *OrderServer {
	return &OrderServer{
		orderService: orderService,
	}
}
