package model

import (
	"gorm.io/gorm"
	"time"
)

type ProductPromo struct {
	gorm.Model     `json:"-"`
	Id             uint64    `json:"id"`
	Title          string    `json:"title"`
	Code           string    `json:"code"`
	PromoTypeId    uint8     `json:"promo_type_id"`
	DiscountAmount uint64    `json:"discount_amount"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	StockLeft      uint64    `json:"stock_left"`
	MerchantId     uint64    `json:"merchant_id"`
}
