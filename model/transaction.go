package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model      `json:"-"`
	Id              uint64 `json:"id"`
	Code            string `json:"code"`
	PaymentMethodId uint64 `json:"payment_method_id"`
	UserId          uint64 `json:"user_id"`
	IsPaid          bool   `json:"is_paid"`
	User            User   `json:"user"`
}
