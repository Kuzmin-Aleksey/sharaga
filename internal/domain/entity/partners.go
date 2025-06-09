package entity

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
