package model

import "gorm.io/gorm"

type UserAddress struct {
	gorm.Model    `json:"-"`
	Id            uint64 `json:"id"`
	Title         string `json:"title"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	UserId        uint64 `json:"user_id"`
	CityId        uint64 `json:"city_id"`
	City          City
}

func (u UserAddress) TableName() string {
	return "user_addresses"
}
