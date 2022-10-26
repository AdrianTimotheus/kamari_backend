package dto

type GetSellerCategoriesRes struct {
	Id            uint64 `json:"id" gorm:"column:id"`
	Title         string `json:"title" gorm:"column:name"`
	ProductsCount uint64 `json:"products_count" gorm:"column:count"`
}
