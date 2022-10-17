package dto

type MerchantVoucherUpdate struct {
	Id             uint64       `json:"id"`
	Title          string       `json:"title"`
	PromoTypeId    uint64       `json:"promo_type_id"`
	DiscountAmount uint64       `json:"discount_amount"`
	MinimumPayment uint64       `json:"minimum_payment"`
	StartDate      JSONNullUnix `json:"start_date"`
	EndDate        JSONNullUnix `json:"end_date"`
	StockLeft      uint         `json:"stock_left"`
}
