package dto

type TransactionPayWithWalletRequest struct {
	Amount uint64 `json:"amount" binding:"required"`
}
