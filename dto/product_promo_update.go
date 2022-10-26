package dto

type promoItemsUpdate struct {
	Id             uint64 `json:"id"`
	ProductId      uint64 `json:"product_id"`
	PromoTypeId    uint8  `json:"promo_type_id"`
	MaximumItem    uint64 `json:"maximum_item"`
	StockLeft      uint64 `json:"stock_left"`
	DiscountAmount uint64 `json:"discount_amount"`
}

type ProductPromoUpdate struct {
	Title      string             `json:"title" binding:"required"`
	EndDate    JSONNullUnix       `json:"end_date" binding:"required"`
	PromoItems []promoItemsUpdate `json:"promo_items" binding:"required"`
}
