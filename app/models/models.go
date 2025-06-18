package models

import "time"

type Error struct {
	Error string `json:"error"`
}

type Id struct {
	Id int `json:"id"`
}

type Tokens struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type User struct {
	Id       int    `json:"id"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Order struct {
	Id        int       `json:"id"`
	CreatorId int       `json:"creator_id"`
	PartnerId int       `json:"partner_id"`
	CreateAt  time.Time `json:"create_at"`
	Price     int       `json:"price"`
}

type OrderProduct struct {
	Id        int `json:"id"`
	OrderId   int `json:"order_id"`
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     int `json:"price"`
}

type OrderProducts struct {
	Order    Order          `json:"order"`
	Products []OrderProduct `json:"products"`
}

type OrderProductInfo struct {
	Order    Order             `json:"order"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	Id          int    `json:"id"`
	Article     int    `json:"article"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinPrice    int    `json:"min_price"`
	SizeX       int    `json:"size_x"`
	SizeY       int    `json:"size_y"`
	SizeZ       int    `json:"size_z"`
	Weight      int    `json:"weight"`
	WeightPack  int    `json:"weight_pack"`
	Quantity    int    `json:"quantity" `
	Price       int    `json:"price"`
}

type Discount struct {
	Discount int `json:"discount"`
}

type Partner struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Director string `json:"director"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	INN      int    `json:"inn"`
	Rating   int    `json:"rating"`
}

type Product struct {
	Id          int    `json:"id"`
	Article     int    `json:"article"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinPrice    int    `json:"min_price"`
	SizeX       int    `json:"size_x"`
	SizeY       int    `json:"size_y"`
	SizeZ       int    `json:"size_z"`
	Weight      int    `json:"weight"`
	WeightPack  int    `json:"weight_pack"`
}

type ProductType struct {
	Id   int     `json:"id" db:"id"`
	Type string  `json:"type" db:"type"`
	K    float64 `json:"k" db:"k"`
}
