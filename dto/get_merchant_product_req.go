package dto

type GetMerchantProductRequest struct {
	SortBy string `form:"sort"`
	Page   int    `form:"page" binding:"required"`
}
