package model

import (
	"gorm.io/gorm"
	"time"
)

type MarketplaceVoucher struct {
	gorm.Model     `json:"-"`
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Code           string     `json:"code"`
	PromoTypeID    int        `json:"promo_type_id"`
	DiscountAmount int        `json:"discount_amount"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	StockLeft      int        `json:"stock_left"`
	PicUrl         string     `json:"pic_url"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
