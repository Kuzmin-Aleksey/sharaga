package entity

type Product struct {
	Id          int    `json:"id" db:"id"`
	Article     int    `json:"article" db:"article"`
	Type        string `json:"type" db:"type"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	MinPrice    int    `json:"min_price" db:"min_price"`
	SizeX       int    `json:"size_x" db:"size_x"`
	SizeY       int    `json:"size_y" db:"size_y"`
	SizeZ       int    `json:"size_z" db:"size_z"`
	Weight      int    `json:"weight" db:"weight"`
	WeightPack  int    `json:"weight_pack" db:"weight_pack"`
}
