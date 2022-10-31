package dto

import (
	"kamari.com/model"
	"time"
)

type OrderMerchantRes struct {
	Id               uint64 `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	PicUrl           string `json:"pic_url"`
	MerchantUsername string `json:"merchant_username"`
}

type UserOrderRes struct {
	Id               uint64                    `json:"id"`
	Invoice          string                    `json:"invoice"`
	UserTotalPayment uint64                    `json:"user_total_payment"`
	TransactionId    uint64                    `json:"transaction_id"`
	Status           string                    `json:"status"`
	Merchant         OrderMerchantRes          `json:"merchant"`
	CourierName      string                    `json:"courier_name"`
	OrderDetail      []model.SellerOrderDetail `json:"order_items" gorm:"-"`
	OrderCreatedAt   time.Time                 `json:"order_created_at,omitempty"`
}
