package dto

import "kamari.com/model"

type ProductPromoItemsResponse struct {
	Id             uint64          `json:"id"`
	DiscountAmount uint64          `json:"discount_amount"`
	MaximumItem    uint64          `json:"maximum_item"`
	StockLeft      uint64          `json:"stock_left"`
	Price          uint64          `json:"price"`
	ProductId      uint64          `json:"product_id"`
	Product        model.Product   `json:"product"`
	ProductPromoId uint64          `json:"product_promo_id"`
	PromoTypeId    uint8           `json:"promo_type_id"`
	PromoType      model.PromoType `json:"promo_type"`
}
