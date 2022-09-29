package model

import "gorm.io/gorm"

type ProductPhoto struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	PicUrl     string `json:"pic_url"`
	ProductId  uint64 `json:"product_id"`
}
