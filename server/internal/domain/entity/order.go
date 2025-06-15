package entity

import "time"

type Order struct {
	Id        int       `json:"id" db:"id"`
	CreatorId int       `json:"creator_id" db:"creator_id"`
	PartnerId int       `json:"partner_id" db:"partner_id"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
}
