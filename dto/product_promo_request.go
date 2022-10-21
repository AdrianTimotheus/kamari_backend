package dto

type promoItems struct {
	ProductPromoId uint64 `json:"product_promo_id"`
	ProductId      uint64 `json:"product_id"`
	PromoTypeId    uint8  `json:"promo_type_id"`
	StockLeft      uint64 `json:"stock_left"`
	MaximumItem    uint64 `json:"maximum_item"`
	DiscountAmount uint64 `json:"discount_amount"`
}

type ProductPromoRequest struct {
	Title      string       `json:"title" binding:"required"`
	StartDate  JSONNullUnix `json:"start_date" binding:"required"`
	EndDate    JSONNullUnix `json:"end_date" binding:"required"`
	PromoItems []promoItems `json:"promo_items" binding:"required"`
}
