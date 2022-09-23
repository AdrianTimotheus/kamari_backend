package service

import (
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type CartService interface {
	GetCartHoverItems(user model.User) ([]dto.CartHoverRes, error)
}

type cartService struct {
	cartRepository repository.CartRepository
	appConfig      config.AppConfig
}

type CSConfig struct {
	CartRepository repository.CartRepository
	AppConfig      config.AppConfig
}

func NewCartService(c CSConfig) CartService {
	return cartService{
		cartRepository: c.CartRepository,
		appConfig:      c.AppConfig,
	}
}

func (c cartService) GetCartHoverItems(user model.User) ([]dto.CartHoverRes, error) {
	result, err := c.cartRepository.GetCartHoverItems(user.Id)
	if err != nil {
		return nil, httperror.InternalServerError("Failed to get cart items")
	}
	return result, nil
}
