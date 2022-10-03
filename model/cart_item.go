package model

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	UserId     uint64 `json:"user_id"`
	VariantId  uint64 `json:"variant_id"`
	Quantity   uint64 `json:"quantity"`
}
