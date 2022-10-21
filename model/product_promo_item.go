package model

import (
	"gorm.io/gorm"
)

type ProductPromoItem struct {
	gorm.Model     `json:"-"`
	Id             uint64    `json:"id"`
	DiscountAmount uint64    `json:"discount_amount"`
	MaximumItem    uint64    `json:"maximum_item"`
	StockLeft      uint64    `json:"stock_left"`
	ProductId      uint64    `json:"product_id"`
	ProductPromoId uint64    `json:"product_promo_id"`
	PromoTypeId    uint8     `json:"promo_type_id"`
	Product        Product   `json:"product"`
	PromoType      PromoType `json:"promo_type"`
}
