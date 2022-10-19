package model

import "gorm.io/gorm"

type Courier struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
}
