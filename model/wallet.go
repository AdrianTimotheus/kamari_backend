package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Balance    uint64 `json:"balance"`
	Pin        string `json:"pin"`
	UserId     uint64 `json:"user_id"`
	User       User   `json:"user"`
}
