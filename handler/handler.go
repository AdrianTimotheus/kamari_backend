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
	wishlistService          service.WishlistService
	sealabsAccountService    service.SealabsAccountService
	merchantVoucherService   service.MerchantVoucherService
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
	WishlistService          service.WishlistService
	SealabsAccountService    service.SealabsAccountService
	MerchantVoucherService   service.MerchantVoucherService
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
		wishlistService:          h.WishlistService,
		sealabsAccountService:    h.SealabsAccountService,
		merchantVoucherService:   h.MerchantVoucherService,
	}
}

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)
