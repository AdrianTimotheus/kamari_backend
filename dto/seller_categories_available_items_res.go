package dto

import "kamari.com/model"

type SellerCategoriesAvailableItemsRes struct {
	Id           uint64                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ThumbnailUrl string                 `json:"thumbnail_url"`
	Variants     []model.ProductVariant `json:"variants" gorm:"foreignKey:product_id;references:id"`
}
