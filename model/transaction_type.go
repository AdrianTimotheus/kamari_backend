package model

import "gorm.io/gorm"

type TransactionType struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Type       string `json:"type"`
}
