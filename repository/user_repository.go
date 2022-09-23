package repository

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/model"
)

type UserRepository interface {
	MatchingCredential(string, password string) (*model.User, error)
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

func (u userRepository) MatchingCredential(email, password string) (*model.User, error) {
	var user model.User
	err := errors.New("Invalid email or password")

	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}

	if checksum := fmt.Sprintf("%X", sha256.Sum256([]byte(password))); checksum != user.Password {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) CreateUser(user *model.User) error {
	return u.db.Create(&user).Error
}
