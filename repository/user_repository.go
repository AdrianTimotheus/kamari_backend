package repository

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
)

type UserRepository interface {
	MatchingCredential(string, password string) (*model.User, error)
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	CheckPassword(id uint64, password string) (*model.User, error)
	GetUserDetail(userId uint64) (*model.User, error)
	EditUserDetail(user model.User) (*model.User, error)
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

func (u userRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := errors.New("Email not found")

	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) CreateUser(user *model.User) error {
	return u.db.Create(&user).Error
}

func (u userRepository) CheckPassword(id uint64, password string) (*model.User, error) {
	var uModel model.User
	err := u.db.Where("id=?", id).
		Where("password=?", password).
		First(&uModel).Error
	if err != nil {
		return nil, err
	}
	return &uModel, nil
}

func (u userRepository) GetUserDetail(userId uint64) (*model.User, error) {
	var user model.User
	err := u.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, httperror.InternalServerError("Failed to get user detail")
	}
	return &user, nil
}

func (u userRepository) EditUserDetail(user model.User) (*model.User, error) {
	var updatedUser model.User
	err := u.db.Where("id = ?", user.Id).Updates(&user).Error
	if err != nil {
		return nil, httperror.InternalServerError("Failed to edit user detail")
	}
	err = u.db.Where("id = ?", user.Id).First(&updatedUser).Error
	if err != nil {
		return nil, httperror.InternalServerError("Failed to get user detail")
	}
	return &updatedUser, nil
}
