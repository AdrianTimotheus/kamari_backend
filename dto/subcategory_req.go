package dto

type NewSubCategoryRequest struct {
	Name       string `json:"name" binding:"required"`
	CategoryId uint64 `json:"categoryId" binding:"required"`
}

type EditSubCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
