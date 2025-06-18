package service

import (
	"app/models"
	"strconv"
)

func (s *Service) NewOrder(order *models.OrderProducts) error {
	id := new(models.Id)

	if err := s.request(id, post, urlOrders, order); err != nil {
		return err
	}

	order.Order.Id = id.Id

	return nil
}

func (s *Service) GetOrders() ([]models.OrderProductInfo, error) {
	orders := make([]models.OrderProductInfo, 0)

	if err := s.request(&orders, get, urlOrders, nil); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) GetOrdersByPartner(partnerId int) ([]models.OrderProductInfo, error) {
	orders := make([]models.OrderProductInfo, 0)

	if err := s.request(&orders, get, urlOrdersByPartner, nil,
		[2]string{"partner_id", strconv.Itoa(partnerId)},
	); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Service) GetPartnerDiscount(partnerId int) (int, error) {
	discount := new(models.Discount)

	if err := s.request(&discount, get, urlOrdersDiscount, nil,
		[2]string{"partner_id", strconv.Itoa(partnerId)},
	); err != nil {
		return 0, err
	}

	return discount.Discount, nil
}
