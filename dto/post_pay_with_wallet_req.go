package dto

type TransactionPayWithWalletRequest struct {
	Amount               uint64         `json:"amount" binding:"required"`
	MarketplaceVoucherID uint64         `json:"marketplace_voucher_id"`
	MerchantVoucherIds   []uint64       `json:"merchant_voucher_ids"`
	Products             []PayedProduct `json:"products" binding:"required"`
}

type PayedProduct struct {
	ProductId uint64 `json:"product_id"`
	VariantId uint64 `json:"variant_id"`
	Quantity  uint64 `json:"quantity"`
}
