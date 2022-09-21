package dto

import (
	"kamari.com/model"
	"time"
)

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DateOfBirth int64  `json:"date_of_birth,omitempty"`
	Phone       string `json:"phone" binding:"required"`
	GenderId    uint8  `json:"gender_id" binding:"required"`
}

func (RegisterRequest) ToModel(req *RegisterRequest) *model.User {
	return &model.User{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		DateOfBirth: time.Unix(req.DateOfBirth, 0),
		Phone:       req.Phone,
		GenderId:    req.GenderId,
		RoleId:      model.RoleBuyer,
	}
}
