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
	walletService      service.WalletService
}

type HConfig struct {
	UserService        service.UserService
	ProvinceService    service.ProvinceService
	MerchantService    service.MerchantService
	CartService        service.CartService
	MarketplaceService service.MarketplaceService
	ProductService     service.ProductService
	WalletService      service.WalletService
}

func New(h HConfig) Handler {
	return Handler{
		userService:        h.UserService,
		provinceService:    h.ProvinceService,
		merchantService:    h.MerchantService,
		cartService:        h.CartService,
		marketplaceService: h.MarketplaceService,
		productService:     h.ProductService,
		walletService:      h.WalletService,
	}
}

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)
