package dto

type SellerCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
