package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model    `json:"-" `
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	PicUrl        string `json:"pic_url"`
	Subcategories []Subcategory
}
