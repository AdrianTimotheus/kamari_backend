package model

import "gorm.io/gorm"

type Province struct {
	gorm.Model `json:"-"`
	Id         int    `json:"id" gorm:"primarykey"`
	Name       string `json:"name"`
}
