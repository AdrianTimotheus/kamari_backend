package dto

type SealabsAccountRequest struct {
	Title      string `json:"title" binding:"required"`
	CardNumber string `json:"card_number" binding:"required,min=8,alphanum"`
	IsDefault  bool   `json:"is_default"`
}
