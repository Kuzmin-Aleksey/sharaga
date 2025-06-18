package service

import "net/http"

const (
	get  = http.MethodGet
	post = http.MethodPost
	put  = http.MethodPut
	del  = http.MethodDelete
)

const (
	urlPing            = "/ping"
	urlLogin           = "/auth/login"
	urlRefresh         = "/auth/refresh"
	urlUserSelf        = "/users/self"
	urlUsers           = "/users"
	urlOrders          = "/orders"
	urlOrdersByPartner = "/orders/by-partner"
	urlOrdersDiscount  = "/orders/discount"
	urlPartners        = "/partners"
	urlProducts        = "/products"
	urlProductsTypes   = "/products/types"
)
