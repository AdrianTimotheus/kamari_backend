package model

import "gorm.io/gorm"

type ReviewPhoto struct {
	gorm.Model      `json:"-"`
	Id              uint64 `json:"id"`
	PicUrl          string `json:"pic_url"`
	ProductReviewId uint64 `json:"product_review_id"`
}
