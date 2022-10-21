package dto

import "kamari.com/model"

type EditProductReq struct {
	Name                    string                 `json:"name"`
	Description             string                 `json:"description"`
	MinimumOrder            uint64                 `json:"minimum_order"`
	MaximumOrder            uint64                 `json:"maximum_order"`
	IsHazardous             bool                   `json:"is_hazardous"`
	Weight                  uint64                 `json:"weight"`
	CategoryId              uint64                 `json:"category_id"`
	SubcategoryId           uint64                 `json:"subcategory_id"`
	ComplementaryCategoryId uint64                 `json:"complementary_category_id"`
	MerchantId              uint64                 `json:"merchant_id"`
	ThumbnailUrl            string                 `json:"thumbnail_url"`
	SellerCategoryId        uint64                 `json:"seller_category_id"`
	ProductVariants         []model.ProductVariant `json:"product_variants"`
}
