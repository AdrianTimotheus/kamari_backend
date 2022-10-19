package model

import "gorm.io/gorm"

type Merchant struct {
	gorm.Model    `json:"-"`
	Id            uint64      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	UserId        uint64      `json:"user_id"`
	UserAddressId uint64      `json:"user_address_id"`
	PicUrl        string      `json:"pic_url"`
	User          User        `json:"user"`
	UserAddress   UserAddress `json:"user_address"`
}
