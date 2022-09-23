package handler

import (
	"errors"
	"kamari.com/service"
)

type Handler struct {
	userService     service.UserService
	provinceService service.ProvinceService
	merchantService service.MerchantService
	cartService     service.CartService
	marketplaceService service.MarketplaceService
}

type HConfig struct {
	UserService     service.UserService
	ProvinceService service.ProvinceService
	MerchantService service.MerchantService
	CartService     service.CartService
	MarketplaceService service.MarketplaceService
}

func New(h HConfig) Handler {
	return Handler{
		userService:     h.UserService,
		provinceService: h.ProvinceService,
		merchantService: h.MerchantService,
		cartService:     h.CartService,
		marketplaceService: c.MarketplaceService,
	}
}

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)
