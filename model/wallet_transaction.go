package model

import "gorm.io/gorm"

type WalletTransaction struct {
	gorm.Model        `json:"-"`
	Id                uint64          `json:"id"`
	SrcWalletId       uint64          `json:"src_wallet_id"`
	DstWalletId       uint64          `json:"dst_wallet_id"`
	Amount            uint64          `json:"amount"`
	Description       string          `json:"description"`
	TransactionTypeId uint64          `json:"transaction_type_id"`
	TransactionType   TransactionType `json:"transaction_type"`
}
