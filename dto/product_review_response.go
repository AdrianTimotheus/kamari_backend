package dto

import (
	"kamari.com/model"
	"time"
)

type reviewer struct {
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
}

func NewReviewer(name, profilePic string) reviewer {
	return reviewer{
		Name:           name,
		ProfilePicture: profilePic,
	}
}

type ProductReviewResponse struct {
	Id                uint64                 `json:"id"`
	Reviewer          reviewer               `json:"reviewer"`
	Score             float32                `json:"score"`
	Feedback          string                 `json:"feedback"`
	SimplifiedVariant string                 `json:"simplified_variant"`
	ProductVariant    []model.ProductVariant `json:"product_variant"`
	ReviewPhotos      []model.ReviewPhoto    `json:"review_photos"`
	CreatedAt         time.Time              `json:"created_at"`
}
