package model

import "gorm.io/gorm"

type CheckoutItem struct {
	gorm.Model       `json:"-"`
	Id               uint64         `json:"id"`
	ProductVariantId uint64         `json:"product_variant_id"`
	MerchantId       uint64         `json:"merchant_id"`
	UserId           uint64         `json:"user_id"`
	Quantity         int            `json:"quantity"`
	ProductVariant   ProductVariant `json:"product_variant,omitempty"`
	Merchant         Merchant       `json:"merchant,omitempty"`
	User             User           `json:"user,omitempty"`
}
