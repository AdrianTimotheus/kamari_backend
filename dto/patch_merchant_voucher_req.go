package dto

type MerchantVoucherUpdate struct {
	Id             uint64       `json:"id"`
	Title          string       `json:"title"`
	PromoTypeId    uint64       `json:"promo_type_id"`
	DiscountAmount uint64       `json:"discount_amount"`
	StartDate      JSONNullUnix `json:"start_date"`
	EndDate        JSONNullUnix `json:"end_date"`
	StockLeft      uint         `json:"stock_left"`
	MerchantId     uint64       `json:"merchant_id,"`
}
