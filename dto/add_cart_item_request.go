package dto

type AddCartItemRequest struct {
	VariantId uint64 `json:"variant_id" binding:"required"`
	Quantity  uint64 `json:"quantity" binding:"required"`
}
