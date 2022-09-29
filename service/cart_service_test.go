package service_test

import (
	"github.com/stretchr/testify/assert"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/service"
	"testing"
	"time"
)

func TestNewCartService(t *testing.T) {
	mockRepo := &mocks.CartRepository{}
	cartService := service.NewCartService(service.CSConfig{
		CartRepository: mockRepo,
	})

	assert.NotNil(t, cartService)
}

func TestCartService_GetCartHoverItems(t *testing.T) {
	t.Run("Should return 5 most recently updated cart items by default", func(t *testing.T) {
		mockRepo := &mocks.CartRepository{}
		cartService := service.NewCartService(service.CSConfig{
			CartRepository: mockRepo,
		})

		data := []dto.CartHoverRes{
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
			dto.CartHoverRes{
				Name:   "T-shirt",
				Price:  100000,
				PicUrl: "https://picsum.photos/id/237/200",
			},
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

		mockRepo.On("GetCartHoverItems", user.Id).Return(data, nil)

		cartData, err := cartService.GetCartHoverItems(user)

		assert.Nil(t, err)
		assert.Equal(t, data, cartData)
	})

	t.Run("Should return error if something went wrong", func(t *testing.T) {
		mockRepo := &mocks.CartRepository{}
		cartService := service.NewCartService(service.CSConfig{
			CartRepository: mockRepo,
		})

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

		mockRepo.On("GetCartHoverItems", user.Id).Return(nil, httperror.InternalServerError("Failed to get cart items"))

		cartData, err := cartService.GetCartHoverItems(user)

		assert.Nil(t, cartData)
		assert.Equal(t, err, httperror.InternalServerError("Failed to get cart items"))
	})
}
