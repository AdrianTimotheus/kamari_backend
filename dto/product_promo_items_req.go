package dto

type PromoItems struct {
	Id             uint64 `json:"id"`
	ProductId      uint64 `json:"product_id"`
	PromoTypeId    uint8  `json:"promo_type_id"`
	MaximumItem    uint64 `json:"maximum_item"`
	DiscountAmount uint64 `json:"discount_amount"`
}
