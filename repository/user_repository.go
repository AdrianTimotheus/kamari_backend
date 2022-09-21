package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type UserRepository interface {
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

type URConfig struct {
	DB *gorm.DB
}

func NewUserRepository(c URConfig) UserRepository {
	return userRepository{
		db: c.DB,
	}
}

func (u userRepository) CreateUser(user *model.User) error {
	return u.db.Create(&user).Error
}
