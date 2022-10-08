package model

import "gorm.io/gorm"

type PromoType struct {
	gorm.Model `json:"-"`
	Id         uint64
	Type       string
}
