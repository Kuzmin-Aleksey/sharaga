package aggregate

import "sharaga/internal/domain/entity"

type OrderProducts struct {
	Order    entity.Order          `json:"order"`
	Products []entity.OrderProduct `json:"products"`
}
