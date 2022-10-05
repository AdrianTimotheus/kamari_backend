package dto

type GetSimilarProductRequest struct {
	Keyword          string `form:"keyword" binding:"required"`
	Category         string `form:"category" binding:"required"`
	CurrentProductId uint64 `form:"CPId" binding:"required"`
}
