package dto

type PayWithSLPReq struct {
	UserID     uint64 `json:"user_id"`
	CardNumber string `json:"card_number" binding:"required"`
	CardID     uint64 `json:"card_id" binding:"required"`
	Amount     int    `json:"amount" binding:"required"`
}
