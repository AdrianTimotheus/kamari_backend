package model

import "time"

type UserOrder struct {
	Id                  uint64              `json:"id"`
	Invoice             string              `json:"invoice"`
	UserTotalPayment    uint64              `json:"user_total_payment"`
	TransactionId       uint64              `json:"transaction_id"`
	Status              string              `json:"status"`
	MerchantId          uint64              `json:"merchant_id"`
	MerchantName        string              `json:"merchant_name"`
	MerchantDescription string              `json:"merchant_description"`
	MerchantPicUrl      string              `json:"merchant_pic_url"`
	MerchantUsername    string              `json:"merchant_username"`
	CourierName         string              `json:"courier_name"`
	OrderDetail         []SellerOrderDetail `json:"order_items" gorm:"-"`
	OrderCreatedAt      time.Time           `json:"order_created_at,omitempty"`
}
