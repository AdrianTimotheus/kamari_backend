package repository

import "gorm.io/gorm"

type OrderStatus struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Status     string `json:"status"`
}
