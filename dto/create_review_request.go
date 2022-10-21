package dto

type CreateReviewRequest struct {
	Score    float32  `json:"score"`
	Feedback string   `json:"feedback"`
	Photos   []string `json:"photos"`
}
