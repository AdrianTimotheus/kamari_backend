package dto

type WalletTransactionDetailResponse struct {
	Source          WalletData `json:"source"`
	Destination     WalletData `json:"destination"`
	Amount          uint64     `json:"amount"`
	TransactionType string     `json:"transaction_type"`
	Description     string     `json:"description"`
	DateOfTime      string     `json:"date_of_time"`
}

type WalletData struct {
	WalletId uint64 `json:"wallet_id"`
	Name     string `json:"name,omitempty"`
}
