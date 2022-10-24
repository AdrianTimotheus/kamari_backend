package model

import "gorm.io/gorm"

type ProductCourier struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	ProductId  uint64 `json:"product_id"`
	CourierId  uint64 `json:"courier_id"`
}
