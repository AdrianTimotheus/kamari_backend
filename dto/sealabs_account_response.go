package dto

import "kamari.com/model"

type SealabsPayAccountResponse struct {
	Id         uint64 `json:"id"`
	Title      string `json:"title"`
	CardNumber string `json:"card_number"`
	IsDefault  bool   `json:"is_default"`
}

func (s *SealabsPayAccountResponse) CreateFromSeaLabsPay(slp *model.SealabsPayAccount) *SealabsPayAccountResponse {
	return &SealabsPayAccountResponse{
		Id:         slp.Id,
		Title:      slp.Title,
		CardNumber: slp.CardNumber,
		IsDefault:  slp.IsDefault,
	}
}
