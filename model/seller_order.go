package model

import "time"

type SellerOrder struct {
	Id               uint64              `json:"id"`
	Invoice          string              `json:"invoice"`
	UserTotalPayment uint64              `json:"user_total_payment"`
	TransactionId    uint64              `json:"transaction_id"`
	MerchantId       uint64              `json:"merchant_id"`
	CourierName      string              `json:"courier_name"`
	WithdrawStatus   bool                `json:"withdraw_status"`
	Status           string              `json:"status"`
	UserUsername     string              `json:"user_username"`
	UserPicUrl       string              `json:"user_pic_url"`
	UserPhone        string              `json:"user_phone"`
	UserEmail        string              `json:"user_email"`
	UserName         string              `json:"user_name"`
	OrderDetail      []SellerOrderDetail `json:"order_items" gorm:"-"`
	OrderCreatedAt   time.Time           `json:"order_created_at,omitempty"`
}
