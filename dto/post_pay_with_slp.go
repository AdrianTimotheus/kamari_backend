package dto

import "kamari.com/model"

type TransactionPaySLPRequest struct {
	Amount               uint64          `json:"amount" binding:"required"`
	CardNumber           string          `json:"card_number" binding:"required"`
	CardID               uint64          `json:"card_id" binding:"required"`
	MarketplaceVoucherID uint64          `json:"marketplace_voucher_id"`
	MerchantVoucherIds   []uint64        `json:"merchant_voucher_ids"`
	Products             []PayedProduct  `json:"products" binding:"required"`
	Orders               []Order         `json:"orders" binding:"required"`
	Items                []ItemToProceed `json:"items" binding:"required"`
	User                 model.User      `json:"user"`
}
