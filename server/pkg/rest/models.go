package rest

type ErrorResponse struct {
	Error string `json:"error"`
}

type IdResponse struct {
	Id int `json:"id"`
}

type DiscountResponse struct {
	Discount int `json:"discount"`
}
