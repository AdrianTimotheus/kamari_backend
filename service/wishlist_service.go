package service

import (
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type WishlistService interface {
	Favorite(user model.User, productId uint64) (*dto.StatusMessageResponse, error)
}

type wishlistService struct {
	wishlistRepository repository.WishlistRepository
}

type WLConfig struct {
	WishlistRepository repository.WishlistRepository
}

func NewWishlistService(c WLConfig) WishlistService {
	return wishlistService{
		wishlistRepository: c.WishlistRepository,
	}
}

func (w wishlistService) Favorite(user model.User, productId uint64) (*dto.StatusMessageResponse, error) {
	response, err := w.wishlistRepository.Favorite(user.Id, productId)
	if err != nil {
		return nil, httperror.InternalServerError("Internal Server Error")
	}
	return response, nil
}
