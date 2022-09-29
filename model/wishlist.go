package model

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	UserId     uint64 `json:"user_id"`
	ProductId  uint64 `json:"product_id"`
}
