package dto

type ValidatePinRequest struct {
	Pin string `json:"pin" binding:"required"`
}
