package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model              `json:"-"`
	Id                      uint64
	Name                    string
	Description             string
	MinimumOrder            uint64
	ThumbnailUrl            string
	IsHazardous             bool
	Weight                  uint64
	ViewCount               uint64
	UnitSold                uint64
	FavoriteCount           uint64
	Availability            bool
	SellerCategoryId        uint64
	CategoryId              uint64
	SubcategoryId           uint64
	MerchantId              uint64
	ComplementaryCategoryId uint64
}
