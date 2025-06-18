package entity

type OrderProduct struct {
	Id        int `json:"id" db:"id"`
	OrderId   int `json:"order_id" db:"order_id"`
	ProductId int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
	Price     int `json:"price" db:"price"`
}
