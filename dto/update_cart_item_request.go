package dto

type UpdateCartItemRequest struct {
	Quantity uint64 `json:"quantity" binding:"required"`
}
