package dto

type MerchantVoucherRequest struct {
	Title          string       `json:"title" binding:"required"`
	Prefix         string       `json:"prefix" binding:"required,min=4"`
	Code           string       `json:"code" binding:"required,min=1,max=5"`
	PromoTypeId    uint64       `json:"promo_type_id" binding:"required"`
	DiscountAmount uint64       `json:"discount_amount" binding:"required"`
	MinimumPayment uint64       `json:"minimum_payment" binding:"required"`
	StartDate      JSONNullUnix `json:"start_date" binding:"required"`
	EndDate        JSONNullUnix `json:"end_date" binding:"required"`
	StockLeft      uint         `json:"stock_left" binding:"required"`
}
