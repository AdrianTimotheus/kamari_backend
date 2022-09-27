package model

import "gorm.io/gorm"

type ProductReview struct {
	gorm.Model    `json:"-"`
	Id            uint64
	OrderDetailId uint64
	Score         float32
	Feedback      string
}
