package dto

type SLPPayReq struct {
	CardNumber   string `json:"card_number"`
	Amount       int    `json:"amount"`
	MerchantCode string `json:"merchant_code"`
	RedirectUrl  string `json:"redirect_url"`
	CallbackUrl  string `json:"callback_url"`
	Signature    string `json:"signature"`
}
