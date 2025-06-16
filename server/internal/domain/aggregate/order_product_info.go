package aggregate

import "sharaga/internal/domain/entity"

type OrderProductInfo struct {
	Order    entity.Order
	Products []ProductQuantity
}
