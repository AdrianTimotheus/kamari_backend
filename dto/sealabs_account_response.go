package dto

type SealabsPayAccountResponse struct {
	Id         uint64 `json:"id"`
	Title      string `json:"title"`
	CardNumber string `json:"card_number"`
	IsDefault  bool   `json:"is_default"`
}
