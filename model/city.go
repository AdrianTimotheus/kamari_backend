package model

import "gorm.io/gorm"

type City struct {
	gorm.Model `json:"-"`
	Id         uint64   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	ProvinceId uint64   `json:"province_id,omitempty"`
	Province   Province `json:"province,omitempty"`
}

func (c City) TableName() string {
	return "cities"
}
