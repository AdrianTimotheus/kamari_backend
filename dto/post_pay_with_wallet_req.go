package dto

type TransactionPayWithWalletRequest struct {
	Amount               uint64 `json:"amount" binding:"required"`
	MerchantVoucherId    uint64 `json:"merchant_voucher_id"`
	MarketplaceVoucherID uint64 `json:"marketplace_voucher_id"`
}
