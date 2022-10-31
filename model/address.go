package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model    `json:"-"`
	Id            uint64 `json:"id"`
	UserId        uint64 `json:"user_id"`
	Title         string `json:"title"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	CityId        uint64 `json:"city_id"`
	DetailAddress string `json:"detail_address"`
	IsDefault     *bool  `json:"is_default"`
}
