package dto

type ReviewResponse struct {
	ReviewCount    int                     `json:"review_count"`
	ReviewScore    float32                 `json:"review_score"`
	ProductReviews []ProductReviewResponse `json:"product_reviews"`
}
