package service_test

import (
	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/service"
	"testing"
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
		req := &dto.RegisterRequest{
			Username: "adbul",
			Password: "somepassword123",
		}
		mockRepo.On("CreateUser", req.ToModel(req)).Return(nil)

		err := s.CreateNewUser(req)

		assert.Nil(t, err)
	})
}
