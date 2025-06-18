package ui

import (
	"app/models"
	"time"
)

// Генерация тестовых данных
func generateTestUsers() []models.User {
	return []models.User{
		{
			Id:       1,
			Role:     models.RoleAdmin,
			Name:     "Администратор",
			Email:    "admin",
			Password: "admin",
		},
		{
			Id:       2,
			Role:     models.RoleManager,
			Name:     "Менеджер",
			Email:    "manager@example.com",
			Password: "manager123",
		},
	}
}

func generateTestPartners() []models.Partner {
	return []models.Partner{
		{
			Id:       1,
			Type:     "ООО",
			Name:     "Технопарк",
			Director: "Иванов И.И.",
			Email:    "tech@example.com",
			Phone:    "+7 999 123-45-67",
			Address:  "Москва, ул. Техническая, 15",
			INN:      1234567890,
			Rating:   5,
		},
		{
			Id:       2,
			Type:     "ИП",
			Name:     "Сервис плюс",
			Director: "Петров П.П.",
			Email:    "service@example.com",
			Phone:    "+7 999 765-43-21",
			Address:  "Санкт-Петербург, пр. Сервисный, 22",
			INN:      9876543210,
			Rating:   4,
		},
	}
}

func generateTestProducts() []models.Product {
	return []models.Product{
		{
			Id:          1,
			Article:     1001,
			Type:        "Электроника",
			Name:        "Ноутбук",
			Description: "Мощный игровой ноутбук",
			MinPrice:    75000,
			SizeX:       35,
			SizeY:       25,
			SizeZ:       3,
			Weight:      2500,
			WeightPack:  3000,
		},
		{
			Id:          2,
			Article:     1002,
			Type:        "Электроника",
			Name:        "Смартфон",
			Description: "Флагманский смартфон",
			MinPrice:    65000,
			SizeX:       15,
			SizeY:       7,
			SizeZ:       1,
			Weight:      200,
			WeightPack:  300,
		},
	}
}

func generateTestOrders() []models.OrderProductInfo {
	now := time.Now()
	return []models.OrderProductInfo{
		{
			Order: models.Order{
				Id:        1,
				CreatorId: 1,
				PartnerId: 1,
				CreateAt:  now.Add(-48 * time.Hour),
				Price:     140000,
			},
			Products: []models.ProductQuantity{
				{
					Id:          1,
					Article:     1001,
					Type:        "Электроника",
					Name:        "Ноутбук",
					Description: "Мощный игровой ноутбук",
					MinPrice:    75000,
					SizeX:       35,
					SizeY:       25,
					SizeZ:       3,
					Weight:      2500,
					WeightPack:  3000,
					Quantity:    1,
				},
				{
					Id:          2,
					Article:     1002,
					Type:        "Электроника",
					Name:        "Смартфон",
					Description: "Флагманский смартфон",
					MinPrice:    65000,
					SizeX:       15,
					SizeY:       7,
					SizeZ:       1,
					Weight:      200,
					WeightPack:  300,
					Quantity:    1,
				},
			},
		},
		{
			Order: models.Order{
				Id:        2,
				CreatorId: 2,
				PartnerId: 2,
				CreateAt:  now.Add(-24 * time.Hour),
				Price:     150000,
			},
			Products: []models.ProductQuantity{
				{
					Id:          1,
					Article:     1001,
					Type:        "Электроника",
					Name:        "Ноутбук",
					Description: "Мощный игровой ноутбук",
					MinPrice:    75000,
					SizeX:       35,
					SizeY:       25,
					SizeZ:       3,
					Weight:      2500,
					WeightPack:  3000,
					Quantity:    2,
				},
			},
		},
	}
}
