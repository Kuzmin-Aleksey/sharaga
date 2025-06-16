package entity

type ProductType struct {
	Id   int     `json:"id" db:"id"`
	Type string  `json:"type_name" db:"type_name"`
	K    float64 `json:"k" db:"k"`
}
