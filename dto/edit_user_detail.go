package dto

import "time"

type EditUserDetailReq struct {
	Name        string    `json:"name,omitempty"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	PicUrl      string    `json:"pic_url,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	GenderId    uint8     `json:"gender_id,omitempty"`
}
