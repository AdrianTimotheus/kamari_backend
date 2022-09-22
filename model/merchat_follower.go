package model

import "gorm.io/gorm"

type MerchantFollower struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	UserId     uint64 `json:"user_id"`
	MerchantId uint64 `json:"merchant_id"`
}
