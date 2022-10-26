package dto

import "kamari.com/model"

type SellerCategoriesItemsRes struct {
	Id               uint64                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	ThumbnailUrl     string                 `json:"thumbnail_url"`
	SellerCategoryId uint64                 `json:"seller_category_id"`
	Variants         []model.ProductVariant `json:"variants" gorm:"foreignKey:product_id;references:id"`
}
