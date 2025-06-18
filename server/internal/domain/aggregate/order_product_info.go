package aggregate

import "sharaga/internal/domain/entity"

type OrderProductInfo struct {
	Order    entity.Order      `json:"order"`
	Products []ProductQuantity `json:"products"`
}
