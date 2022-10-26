package dto

type ComplementaryCategoryResponse struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	SubCategoryId uint64 `json:"subcategory_id"`
}
