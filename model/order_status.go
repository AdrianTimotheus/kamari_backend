package model

import "gorm.io/gorm"

type OrderStatus struct {
	gorm.Model `json:"-"`
	Id         uint64
	Status     string
}
