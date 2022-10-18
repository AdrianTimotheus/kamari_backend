package dto

type SLPPayRes struct {
	RedirectUrl string `json:"redirectUrl"`
}

func (slp *SLPPayRes) Create(url string) *SLPPayRes {
	return &SLPPayRes{RedirectUrl: url}
}
