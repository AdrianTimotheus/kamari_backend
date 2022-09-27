package model

import "gorm.io/gorm"

type Subcategory struct {
	gorm.Model `json:"-" `
	Id         uint64
	Name       string
	CategoryId uint64
}
