package aggregate

import "sharaga/internal/domain/entity"

type ProductWithType struct {
	Product entity.Product   `json:"product"`
	Type    *ProductWithType `json:"type"`
}
