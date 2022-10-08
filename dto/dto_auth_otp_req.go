package dto

type AuthOTPRequest struct {
	Code string `json:"code" binding:"required"`
}
