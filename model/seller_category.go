package model

import "gorm.io/gorm"

type SellerCategory struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	MerchantId uint64 `json:"merchantId"`
}
