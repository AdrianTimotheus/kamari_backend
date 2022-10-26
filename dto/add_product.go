package dto

import "kamari.com/model"

type AddProductReq struct {
	Name                    string                 `json:"name" binding:"required"`
	Description             string                 `json:"description" binding:"required"`
	MinimumOrder            uint64                 `json:"minimum_order" binding:"required,gte=1"`
	MaximumOrder            uint64                 `json:"maximum_order" binding:"required,gte=1"`
	IsHazardous             *bool                  `json:"is_hazardous" binding:"required"`
	Weight                  uint64                 `json:"weight" binding:"required"`
	CategoryId              uint64                 `json:"category_id" binding:"required"`
	SubcategoryId           uint64                 `json:"subcategory_id" binding:"required"`
	ComplementaryCategoryId uint64                 `json:"complementary_category_id"`
	ThumbnailUrl            string                 `json:"thumbnail_url" binding:"required"`
	Photos                  []string               `json:"photos" binding:"required"`
	ProductVariants         []model.ProductVariant `json:"product_variants" binding:"required"`
}
