package httpserver

type productService interface {
}

type ProductServer struct {
	productService productService
}

func NewProductServer(productService productService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}
