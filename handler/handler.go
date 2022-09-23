package handler

import (
	"errors"
	"kamari.com/service"
)

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)

type Handler struct {
	userService        service.UserService
	provinceService    service.ProvinceService
	merchantService    service.MerchantService
	marketplaceService service.MarketplaceService
}

type HConfig struct {
	UserService        service.UserService
	ProvinceService    service.ProvinceService
	MerchantService    service.MerchantService
	MarketplaceService service.MarketplaceService
}

func New(c HConfig) Handler {
	return Handler{
		userService:        c.UserService,
		provinceService:    c.ProvinceService,
		merchantService:    c.MerchantService,
		marketplaceService: c.MarketplaceService,
	}
}
