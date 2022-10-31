package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model               `json:"-"`
	Id                       uint64
	Invoice                  string
	UserTotalPayment         uint64
	MerchantIncome           uint64
	MerchantVoucherId        *uint64 `gorm:"default:null"`
	MarketplaceVoucherId     *uint64 `gorm:"default:null"`
	DstAddressId             uint64
	TransactionId            uint64
	OrderStatusId            uint64
	OrderStatus              OrderStatus
	WithdrawStatus           bool
	MerchantId               uint64
	Merchant                 Merchant
	CourierId                uint64
	UserId                   uint64
	ShippingFee              uint64
	TotalProductPrice        uint64
	MerchantTotalDiscount    uint64
	MarketplaceTotalDiscount uint64
	Transaction              Transaction
	Courier                  Courier
	User                     User
}
