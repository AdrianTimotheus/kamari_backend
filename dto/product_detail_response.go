package dto

import "kamari.com/model"

type ProductDetailResponse struct {
	Id                      uint64                 `json:"id,omitempty"`
	Title                   string                 `json:"title,omitempty"`
	Description             string                 `json:"description,omitempty"`
	MinimumOrder            uint64                 `json:"minimum-order,omitempty"`
	ThumbnailUrl            string                 `json:"thumbnail-url,omitempty"`
	IsHazardous             bool                   `json:"is-hazardous,omitempty"`
	Weight                  uint64                 `json:"weight,omitempty"`
	ViewCount               uint64                 `json:"view-count,omitempty"`
	UnitSold                uint64                 `json:"unit-sold,omitempty"`
	FavouriteCount          uint64                 `json:"favourite-count,omitempty"`
	Stock                   int                    `json:"stock,omitempty"`
	Availability            bool                   `json:"availability,omitempty"`
	SellerCategoryId        uint64                 `json:"seller-category-id,omitempty"`
	CategoryId              uint64                 `json:"category-id,omitempty"`
	SubcategoryId           uint64                 `json:"subcategory-id,omitempty"`
	ComplementaryCategoryId uint64                 `json:"complementary-category-id,omitempty"`
	MerchantID              uint64                 `json:"merchant-id,omitempty"`
	ProductPhotos           []model.ProductPhoto   `json:"product_photos"`
	ProductVariants         []model.ProductVariant `json:"product_variants"`
}
