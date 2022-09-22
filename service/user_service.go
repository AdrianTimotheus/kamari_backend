package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
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

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(req.Password)))
	user.Password = hash

	err := u.userRepository.CreateUser(user)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr):
			return httperror.InternalServerError("duplicate key values for unique constraint")
		default:
			return httperror.InternalServerError(err.Error())
		}
	}

	return nil
}
