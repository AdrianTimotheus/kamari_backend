package dto

type SubCategoryResponse struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	CategoryId uint64 `json:"category_id"`
}
