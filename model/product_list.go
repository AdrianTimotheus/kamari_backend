package model

import (
	"gorm.io/gorm"
)

type ProductList struct {
	gorm.Model       `json:"-"`
	Id               uint64
	Name             string
	ThumbnailUrl     string
	City             string
	Price            uint64
	FavouriteCount   uint64
	UnitSold         uint64
	Availability     bool
	CategoryId       uint64
	MerchantId       uint64
	Seller           string
	SellerCategoryId uint64
	Merchant         Merchant
	Category         Category
	SellerCategory   SellerCategory
}

func (p ProductList) TableName() string {
	return "products"
}
