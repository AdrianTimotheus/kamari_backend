package model

import "gorm.io/gorm"

type MerchantVoucher struct {
	gorm.Model  `json:"-"`
	Id          uint64
	Title       string
	Code        string
	PromoTypeId uint64
	//DiscountAmount
}
