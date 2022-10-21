package model

import (
	"gorm.io/gorm"
	"time"
)

type ProductPromo struct {
	gorm.Model        `json:"-"`
	Id                uint64             `json:"id"`
	Title             string             `json:"title"`
	StartDate         time.Time          `json:"start_date"`
	EndDate           time.Time          `json:"end_date"`
	MerchantId        uint64             `json:"merchant_id"`
	ProductPromoItems []ProductPromoItem `json:"product_promo_items"`
}
