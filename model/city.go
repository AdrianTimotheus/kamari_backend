package model

import "gorm.io/gorm"

type City struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
}

func (c City) TableName() string {
	return "cities"
}
