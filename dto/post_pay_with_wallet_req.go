package dto

type TransactionPayWithWalletRequest struct {
	Amount               uint64   `json:"amount" binding:"required"`
	MarketplaceVoucherID uint64   `json:"marketplace_voucher_id"`
	MerchantVoucherIds   []uint64 `json:"merchant_voucher_ids"`
}
