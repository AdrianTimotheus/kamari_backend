package model

import (
	"gorm.io/gorm"
)

type ProductList struct {
	gorm.Model     `json:"-"`
	Id             uint64
	Name           string
	ThumbnailUrl   string
	City           string
	Price          uint64
	FavouriteCount uint64
	UnitSold       uint64
	Availability   bool
	CategoryId     uint64
	MerchantId     uint64
	Seller         string
	Merchant       Merchant
	Category       Category
}

func (p ProductList) TableName() string {
	return "products"
}
