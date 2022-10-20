package model

import "gorm.io/gorm"

type WalletPayment struct {
	gorm.Model    `json:"-"`
	Id            uint64      `json:"id"`
	WalletId      uint64      `json:"wallet_id"`
	TransactionId uint64      `json:"transaction_id"`
	Wallet        Wallet      `json:"wallet,omitempty"`
	Transaction   Transaction `json:"transaction,omitempty"`
}
