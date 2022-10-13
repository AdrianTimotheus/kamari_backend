package model

import (
	"gorm.io/gorm"
	"time"
)

type MerchantVoucher struct {
	gorm.Model     `json:"-"`
	Id             uint64
	Title          string
	Code           string
	PromoTypeId    uint64
	DiscountAmount uint64
	StartDate      time.Time
	EndDate        time.Time
	StockLeft      uint
	MerchantId     uint64
}
