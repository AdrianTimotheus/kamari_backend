package handler

import (
	"errors"
	"kamari.com/service"
)

type Handler struct {
	userService        service.UserService
	provinceService    service.ProvinceService
	merchantService    service.MerchantService
	cartService        service.CartService
	marketplaceService service.MarketplaceService
	productService     service.ProductService
}

type HConfig struct {
	UserService        service.UserService
	ProvinceService    service.ProvinceService
	MerchantService    service.MerchantService
	CartService        service.CartService
	MarketplaceService service.MarketplaceService
	ProductService     service.ProductService
}

func New(h HConfig) Handler {
	return Handler{
		userService:        h.UserService,
		provinceService:    h.ProvinceService,
		merchantService:    h.MerchantService,
		cartService:        h.CartService,
		marketplaceService: h.MarketplaceService,
		productService:     h.ProductService,
	}
}

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)
