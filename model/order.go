package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model           `json:"-"`
	Id                   uint64
	Invoice              string
	UserTotalPayment     uint64
	MerchantIncome       uint64
	ShippingFee          uint64
	TotalProductsPrice   uint64
	MerchantVoucherId    uint64
	MarketplaceVoucherId uint64
	DstAddressId         uint64
	TransactionId        uint64
	Transaction          Transaction
	OrderStatusId        uint64
	WithdrawStatus       bool
	MerchantId           uint64
	CourierId            uint64
	UserId               uint64
}
