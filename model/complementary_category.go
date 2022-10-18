package model

import "gorm.io/gorm"

type Complementary struct {
	gorm.Model    `json:"-" `
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	SubcategoryId uint64 `json:"subcategory_id"`
}
