package dto

type GetMerchantProductRequest struct {
	SortBy string `form:"sort"`
	Page   int    `form:"page" binding:"required"`
	Limit  int    `form:"limit" binding:"required"`
	Search string `form:"search"`
}
