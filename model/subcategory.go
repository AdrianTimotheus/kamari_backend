package model

import "gorm.io/gorm"

type Subcategory struct {
	gorm.Model `json:"-" `
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	CategoryId uint64 `json:"category_id"`
}
