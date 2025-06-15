package entity

type Partner struct {
	Id       int    `json:"id" db:"id"`
	Type     string `json:"type" db:"type"`
	Name     string `json:"name" db:"name"`
	Director string `json:"director" db:"director"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Address  string `json:"address" db:"address"`
	INN      int    `json:"inn" db:"inn"`
	Rating   int    `json:"rating" db:"rating"`
}
