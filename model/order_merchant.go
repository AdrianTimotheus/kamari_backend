package model

import "gorm.io/gorm"

type OrderMerchant struct {
	gorm.Model           `json:"-"`
	Id                   uint64
	Invoice              string
	UserTotalPayment     uint64
	MerchantIncome       uint64
	ShippingFee          uint64
	TotalProductsPrice   uint64
	MerchantVoucherId    uint64
	MerchantVoucher      MerchantVoucher
	MarketplaceVoucherId uint64
	MarketplaceVoucher   MarketplaceVoucher
	DstAddressId         uint64
	UserAddress          UserAddress
	TransactionId        uint64
	Transaction          Transaction
	OrderStatusId        uint64
	//OrderStatus          OrderStatus // todo: create orderstatus model
	WithdrawStatus bool
	MerchantId     uint64
	CourierId      uint64
	UserId         uint64
}
