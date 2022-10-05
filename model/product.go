package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model              `json:"-"`
	Id                      uint64 `json:"id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	MinimumOrder            uint64 `json:"minimum_order"`
	MaximumOrder            uint64 `json:"maximum_order"`
	ThumbnailUrl            string
	IsHazardous             bool
	Weight                  uint64
	ViewCount               uint64
	UnitSold                uint64
	FavouriteCount          uint64
	Availability            bool
	SellerCategoryId        uint64
	CategoryId              uint64
	SubcategoryId           uint64
	ComplementaryCategoryId uint64
	MerchantID              uint64
	Merchant                Merchant
	Category                Category
	Subcategory             Subcategory
	ProductPhotos           []ProductPhoto
	ProductVariants         []ProductVariant
}
