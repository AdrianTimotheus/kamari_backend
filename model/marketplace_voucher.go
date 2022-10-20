package model

import (
	"gorm.io/gorm"
	"time"
)

type MarketplaceVoucher struct {
	gorm.Model     `json:"-"`
	ID             uint64    `json:"id"`
	Title          string    `json:"title"`
	Code           string    `json:"code"`
	PromoTypeID    uint64    `json:"promo_type_id"`
	DiscountAmount uint64    `json:"discount_amount"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	StockLeft      int       `json:"stock_left"`
	PicUrl         string    `json:"pic_url"`
	PromoType      PromoType `json:"promo_type,omitempty"`
}
