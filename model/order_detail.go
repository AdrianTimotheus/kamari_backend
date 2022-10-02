package model

import "gorm.io/gorm"

type OrderDetail struct {
	gorm.Model       `json:"-"`
	Id               uint64
	Quantity         uint64
	OrderId          uint64
	Order            Order
	ProductVariantId uint64
	CourierRouteId   uint64
	ProductReview    ProductReview
}
