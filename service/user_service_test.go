package service_test

import (
	"crypto/sha256"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/service"
	"testing"
	"time"
)

func TestNewUserService(t *testing.T) {
	mockRepo := &mocks.UserRepository{}
	userService := service.NewUserService(service.USConfig{
		UserRepository: mockRepo,
	})

	assert.NotNil(t, userService)
}

func Test_userService_CreateNewUser(t *testing.T) {
	t.Run("should thrown error if password contains username", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		s := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		req := &dto.RegisterRequest{
			Username: "adbul",
			Password: "adbul",
		}

		err := s.CreateNewUser(req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, httperror.BadRequestError("password cannot contains username", "INVALID_PASSWORD_PAYLOAD"))
	})

	t.Run("should return nil when user creation is running smoothly", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		s := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		p := fmt.Sprintf("%X", sha256.Sum256([]byte("password123")))
		req := &dto.RegisterRequest{
			Username: "adbul",
			Password: p,
		}
		mockRepo.On("CreateUser", req.ToModel(req)).Return(nil)

		err := s.CreateNewUser(&dto.RegisterRequest{
			Username: "adbul",
			Password: "password123",
		})

		assert.Nil(t, err)
	})

	t.Run("should thrown an error when trying to register with registered email/username", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		s := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		p := fmt.Sprintf("%X", sha256.Sum256([]byte("password123")))
		req := &dto.RegisterRequest{
			Username: "adbul",
			Password: p,
		}
		mockRepo.On("CreateUser", req.ToModel(req)).Return(&pgconn.PgError{})

		err := s.CreateNewUser(&dto.RegisterRequest{
			Username: "adbul",
			Password: "password123",
		})

		assert.Equal(t, httperror.InternalServerError("duplicate key values for unique constraint"), err)
		assert.Error(t, err)
	})

	t.Run("should thrown an error when trying to register with password length less than 8 chars", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		s := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		p := fmt.Sprintf("%X", sha256.Sum256([]byte("123")))
		req := &dto.RegisterRequest{
			Username: "adbul",
			Password: p,
		}
		mockRepo.On("CreateUser", req.ToModel(req)).Return(httperror.InternalServerError("failed binding validation"))

		err := s.CreateNewUser(&dto.RegisterRequest{
			Username: "adbul",
			Password: "123",
		})

		assert.Error(t, err)
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("Should return error when email and password don't match", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		userService := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		//p := fmt.Sprintf("%X", sha256.Sum256([]byte("password")))
		req := &dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "password",
		}
		mockRepo.On("MatchingCredential", req.Email, req.Password).Return(nil, httperror.InvalidLogin("Invalid email or password"))

		tokenResponse, err := userService.Login(&dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "password",
		})

		assert.NotNil(t, err)
		assert.Nil(t, tokenResponse)
	})

	t.Run("Should return token if email and password match", func(t *testing.T) {
		mockRepo := &mocks.UserRepository{}
		userService := service.NewUserService(service.USConfig{
			UserRepository: mockRepo,
		})
		//p := fmt.Sprintf("%X", sha256.Sum256([]byte("password")))
		req := &dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "password",
		}

		user := model.User{
			Id:          4,
			Name:        "Adrian",
			Username:    "adrian",
			Email:       "adrian@gmail.com",
			Phone:       "08123456789",
			DateOfBirth: time.Unix(0, 0).UTC(),
			GenderId:    1,
			RoleId:      2,
		}

		mockRepo.On("MatchingCredential", req.Email, req.Password).Return(&user, nil)

		tokenResponse, err := userService.Login(&dto.LoginRequest{
			Email:    "adrian@gmail.com",
			Password: "password",
		})

		assert.Nil(t, err)
		assert.NotNil(t, tokenResponse)
	})
}
