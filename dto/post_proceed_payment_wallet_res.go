package dto

type ProcessPaymentWithWalletResponse struct {
	Status        bool   `json:"status"`
	TransactionId uint64 `json:"transaction_id"`
}
