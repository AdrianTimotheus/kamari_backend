package model

import "gorm.io/gorm"

type OtpCode struct {
	gorm.Model `json:"-"`
	Id         uint64 `json:"id"`
	Code       string `json:"code"`
	UserId     uint64 `json:"user_id"`
	User       User   `json:"user"`
}
