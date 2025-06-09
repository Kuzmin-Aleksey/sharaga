package entity

type User struct {
	Id       int    `json:"id" db:"id"`
	Role     string `json:"role" db:"role"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password,omitempty" db:"password"`
}
