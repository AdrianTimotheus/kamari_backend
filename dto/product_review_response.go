package dto

import (
	"kamari.com/model"
	"time"
)

type ProductReviewResponse struct {
	Id                uint64                 `json:"id"`
	Score             float32                `json:"score"`
	Feedback          string                 `json:"feedback"`
	SimplifiedVariant string                 `json:"simplified_variant"`
	ProductVariant    []model.ProductVariant `json:"product_variant"`
	ReviewPhotos      []model.ReviewPhoto    `json:"review_photos"`
	CreatedAt         time.Time              `json:"created_at"`
}
