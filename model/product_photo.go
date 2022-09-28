package model

import "gorm.io/gorm"

type ProductPhoto struct {
	gorm.Model `json:"-"`
	Id         uint64
	PicUrl     string
	ProductId  uint64
}
