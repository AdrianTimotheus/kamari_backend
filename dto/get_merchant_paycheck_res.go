package dto

type MerchantPaycheckResponse struct {
	Merchant MerchantDetailResponse `json:"merchant"`
	Income   uint64                 `json:"income"`
}
