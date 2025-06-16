package entity

type ProductType struct {
	Id   int     `json:"id" db:"id"`
	Type string  `json:"type" db:"type"`
	K    float64 `json:"k" db:"k"`
}
