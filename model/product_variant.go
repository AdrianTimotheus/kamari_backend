package model

import "gorm.io/gorm"

type ProductVariant struct {
	gorm.Model      `json:"-"`
	ID              uint64 `json:"id"`
	MainVariant     string `json:"main_variant"`
	MainVariantItem string `json:"main_variant_item"`
	SubVariant      string `json:"sub_variant"`
	SubVariantItem  string `json:"sub_variant_item"`
	Stock           int    `json:"stock"`
	ProductStatus   bool   `json:"product_status"`
	Price           uint64 `json:"price"`
	PicUrl          string `json:"pic_url"`
	ProductPromoID  int    `json:"product_promo_id"`
	ProductID       int    `json:"product_id"`
}
