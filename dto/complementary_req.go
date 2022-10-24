package dto

type NewComplementaryCategoryRequest struct {
	Name          string `json:"name" binding:"required"`
	SubCategoryId uint64 `json:"subCategoryId" binding:"required"`
}

type EditComplementaryCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
