package handler

import (
	"errors"

	"kamari.com/service"
)

type Handler struct {
	userService              service.UserService
	provinceService          service.ProvinceService
	merchantService          service.MerchantService
	cartService              service.CartService
	marketplaceService       service.MarketplaceService
	productService           service.ProductService
	walletService            service.WalletService
	categoryService          service.CategoryService
	productReviewService     service.ProductReviewService
	walletTransactionService service.WalletTransactionService
}

type HConfig struct {
	UserService              service.UserService
	ProvinceService          service.ProvinceService
	MerchantService          service.MerchantService
	CartService              service.CartService
	MarketplaceService       service.MarketplaceService
	ProductService           service.ProductService
	WalletService            service.WalletService
	CategoryService          service.CategoryService
	ProductReviewService     service.ProductReviewService
	WalletTransactionService service.WalletTransactionService
}

func New(h HConfig) Handler {
	return Handler{
		userService:              h.UserService,
		provinceService:          h.ProvinceService,
		merchantService:          h.MerchantService,
		cartService:              h.CartService,
		marketplaceService:       h.MarketplaceService,
		productService:           h.ProductService,
		walletService:            h.WalletService,
		categoryService:          h.CategoryService,
		productReviewService:     h.ProductReviewService,
		walletTransactionService: h.WalletTransactionService,
	}
}

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)
