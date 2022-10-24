package model

import "gorm.io/gorm"

type MerchantCourier struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	CourierId  uint64 `json:"courier_id"`
	MerchantId uint64 `json:"merchant_id"`
}
