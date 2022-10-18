package dto

type PostCheckoutRequest struct {
	ProductVariantId uint64 `json:"product_variant_id" binding:"required"`
	Quantity         int    `json:"quantity" binding:"required"`
	MerchantId       uint64 `json:"merchant_id" binding:"required"`
}
