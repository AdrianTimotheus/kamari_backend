package dto

type WalletTransactionResponse struct {
	Page      int                 `json:"page"`
	TotalPage int                 `json:"total_page"`
	UserId    uint64              `json:"user_id"`
	WalletID  uint64              `json:"wallet_id"`
	Records   []TransactionRecord `json:"records"`
}

type TransactionRecord struct {
	Id                uint64 `json:"id"`
	SourceWallet      uint64 `json:"source_wallet"`
	DestinationWallet uint64 `json:"destination_wallet"`
	Amount            uint64 `json:"amount"`
	TransactionType   string `json:"transaction_type"`
	Description       string `json:"description"`
	DateOfTime        string `json:"date_of_time"`
}
