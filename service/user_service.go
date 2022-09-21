package service

import (
	"golang.org/x/crypto/bcrypt"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type UserService interface {
	CreateNewUser(req *dto.RegisterRequest) error
}

type userService struct {
	userRepository repository.UserRepository
}

type USConfig struct {
	UserRepository repository.UserRepository
}

func NewUserService(c USConfig) UserService {
	return userService{
		userRepository: c.UserRepository,
	}
}

func (u userService) CreateNewUser(req *dto.RegisterRequest) error {
	if helper.EmailNotValid(req.Username, req.Password) {
		return httperror.BadRequestError("password cannot contains username", "INVALID_PASSWORD_PAYLOAD")
	}

	user := req.ToModel(req)

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	err := u.userRepository.CreateUser(user)
	if err != nil {
		return httperror.InternalServerError(err.Error())
	}

	return nil
}
