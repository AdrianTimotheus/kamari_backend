package dto

import (
	"kamari.com/model"
	"time"
)

type MerchantVoucherResponse struct {
	Id             uint64          `json:"id"`
	Title          string          `json:"title"`
	Code           string          `json:"code"`
	StockLeft      uint            `json:"stock_left"`
	DiscountAmount uint64          `json:"discount_amount"`
	MinimumPayment uint64          `json:"minimum_payment,omitempty"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        time.Time       `json:"end_date"`
	MerchantId     uint64          `json:"merchant_id,omitempty"`
	PromoTypeId    uint64          `json:"promo_type_id"`
	PromoType      model.PromoType `json:"promo_type"`
}
