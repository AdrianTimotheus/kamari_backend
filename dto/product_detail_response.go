package dto

import "kamari.com/model"

type ProductDetailResponse struct {
	Id                      uint64                 `json:"id"`
	Title                   string                 `json:"title"`
	Description             string                 `json:"description"`
	MinimumOrder            uint64                 `json:"minimum_order"`
	MaximumOrder            uint64                 `json:"maximum_order"`
	ThumbnailUrl            string                 `json:"thumbnail_url"`
	IsHazardous             bool                   `json:"is_hazardous"`
	Weight                  uint64                 `json:"weight"`
	ViewCount               uint64                 `json:"view_count"`
	UnitSold                uint64                 `json:"unit_sold"`
	FavouriteCount          uint64                 `json:"favourite_count"`
	Stock                   int                    `json:"stock"`
	Availability            bool                   `json:"availability"`
	SellerCategoryId        uint64                 `json:"seller_category_id"`
	CategoryId              uint64                 `json:"category_id"`
	SubcategoryId           uint64                 `json:"subcategory_id,omitempty"`
	ComplementaryCategoryId uint64                 `json:"complementary_category_id"`
	MerchantID              uint64                 `json:"merchant_id"`
	MerchantUsername        string                 `json:"merchant_username"`
	Category                model.Category         `json:"category"`
	Subcategory             model.Subcategory      `json:"subcategory"`
	ProductPhotos           []model.ProductPhoto   `json:"product_photos"`
	ProductVariants         []model.ProductVariant `json:"product_variants"`
	PriceRange              []uint64               `json:"price_range"`
	ReviewCount             int                    `json:"review_count"`
	ReviewScore             float32                `json:"review_score"`
	Promo                   *PromoItems            `json:"promo"`
}
