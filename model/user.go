package model

import (
	"gorm.io/gorm"
	"time"
)

type RoleId uint

const (
	_ RoleId = iota
	RoleAdmin
	RoleBuyer
	RoleSeller
)

type User struct {
	gorm.Model  `json:"-"`
	Id          uint64    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Password    string    `json:"-,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	DateOfBirth time.Time `json:"date_of_birth,omitempty"`
	GenderId    uint8     `json:"gender_id,omitempty"`
	CityId      uint64    `json:"city_id,omitempty"`
	RoleId      RoleId    `json:"role_id,omitempty"`
}
