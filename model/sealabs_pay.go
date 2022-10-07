package model

import "gorm.io/gorm"

type SealabsPayAccount struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Title      string
	CardNumber string
	UserId     uint64
	IsDefault  bool
}
