package dto

type WalletTransactionRequest struct {
	Page int `form:"page" binding:"required"`
}
