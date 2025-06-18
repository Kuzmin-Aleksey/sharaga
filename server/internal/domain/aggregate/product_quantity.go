package aggregate

import "sharaga/internal/domain/entity"

type ProductQuantity struct {
	entity.Product
	Quantity int `json:"quantity" db:"quantity"`
	Price    int `json:"price" db:"price"`
}
