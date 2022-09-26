package model

import "gorm.io/gorm"

type City struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	ProvinceId uint64 `json:"province_id"`
	Province   Province
}

func (c City) TableName() string {
	return "cities"
}
