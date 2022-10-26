package dto

import (
	"time"
)

type ProductPromoResponse struct {
	Id                uint64                      `json:"id"`
	Title             string                      `json:"title"`
	StartDate         time.Time                   `json:"start_date"`
	EndDate           time.Time                   `json:"end_date"`
	MerchantId        uint64                      `json:"merchant_id"`
	HasEnded          bool                        `json:"has_ended"`
	ProductPromoItems []ProductPromoItemsResponse `json:"product_promo_items"`
}
