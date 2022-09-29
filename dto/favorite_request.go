package dto

type FavoriteRequest struct {
	ProductId uint64 `json:"product_id" binding:"required"`
}
