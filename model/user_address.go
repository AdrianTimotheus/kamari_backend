package model

import "gorm.io/gorm"

type UserAddress struct {
	gorm.Model    `json:"-"`
	Id            uint64 `json:"id"`
	UserId        uint64 `json:"user_id"`
	User          User   `json:"user,omitempty"`
	Title         string `json:"title"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	CityId        uint64 `json:"city_id"`
	City          City   `json:"city,omitempty"`
	DetailAddress string `json:"detail_address"`
	IsDefault     bool   `json:"is_default" gorm:"column:is_default"`
}

func (u UserAddress) TableName() string {
	return "user_addresses"
}
