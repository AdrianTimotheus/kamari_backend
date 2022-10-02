package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model    `json:"-"`
	Id            uint64
	Code          string
	PaymentMethod uint64
	UserId        uint64
	User          User
}
