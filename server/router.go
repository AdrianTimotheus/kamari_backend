package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kamari.com/handler"
	"kamari.com/middleware"
	"kamari.com/service"
)

type RouterConfig struct {
	UserService               service.UserService
	ProvinceService           service.ProvinceService
	MerchantService           service.MerchantService
	CartService               service.CartService
	MarketplaceService        service.MarketplaceService
	ProductService            service.ProductService
	WalletService             service.WalletService
	CategoryService           service.CategoryService
	ProductReviewService      service.ProductReviewService
	WalletTransactionService  service.WalletTransactionService
	WishlistService           service.WishlistService
	SealabsAccountService     service.SealabsAccountService
	MerchantVoucherService    service.MerchantVoucherService
	AddressService            service.AddressService
	CityService               service.CityService
	CheckoutItemService       service.CheckoutItemService
	CourierService            service.CourierService
	MarketplaceVoucherService service.MarketplaceVoucherService
	TransactionService        service.TransactionService
	ProductPromoService       service.ProductPromoService
	OrderService              service.OrderService
}

func NewRouter(c RouterConfig) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.New(handler.HConfig{
		UserService:               c.UserService,
		ProvinceService:           c.ProvinceService,
		MerchantService:           c.MerchantService,
		ProductService:            c.ProductService,
		CartService:               c.CartService,
		MarketplaceService:        c.MarketplaceService,
		WalletService:             c.WalletService,
		CategoryService:           c.CategoryService,
		ProductReviewService:      c.ProductReviewService,
		WalletTransactionService:  c.WalletTransactionService,
		WishlistService:           c.WishlistService,
		SealabsAccountService:     c.SealabsAccountService,
		MerchantVoucherService:    c.MerchantVoucherService,
		AddressService:            c.AddressService,
		CityService:               c.CityService,
		CheckoutItemService:       c.CheckoutItemService,
		CourierService:            c.CourierService,
		MarketplaceVoucherService: c.MarketplaceVoucherService,
		TransactionService:        c.TransactionService,
		ProductPromoService:       c.ProductPromoService,
		OrderService:              c.OrderService,
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("route not found, go to %v/docs for documentation\n", os.Getenv("BASE_URL")),
		})
	})

	api := router.Group("/api")
	{
		// auth
		api.POST("/users", h.Register, middleware.ErrorHandler)
		api.POST("/login", h.Login, middleware.ErrorHandler)
		api.POST("/login-with-google", h.LoginWithGoogle, middleware.ErrorHandler)
		api.GET("/refresh-token", h.GetNewTokenFromRefresh, middleware.ErrorHandler)
		api.GET("/users", middleware.AuthorizeJWT, h.GetUserDetail, middleware.ErrorHandler)
		api.GET("/users/seller", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetUserDetail, middleware.ErrorHandler)
		api.GET("/users/admin", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.GetUserDetail, middleware.ErrorHandler)
		api.GET("/logout", h.Logout)

		// guest
		api.Static("/docs", "swaggerui")
		api.GET("/provinces", h.GetAllProvinces, middleware.ErrorHandler)
		api.GET("/sellers/:username", h.GetMerchantDetail, middleware.ErrorHandler)
		api.GET("/sellers/:username/products", h.GetSellerProduct, middleware.ErrorHandler)
		api.GET("/marketplace/promo/banners", h.GetMarketplacePromoBanner, middleware.ErrorHandler)
		api.GET("/products", h.GetProducts, middleware.ErrorHandler)
		api.GET("/products/similar", h.GetSimilarProducts, middleware.ErrorHandler)
		api.GET("/products/recommend", h.GetRecommendedProduct, middleware.ErrorHandler)
		api.GET("/products/:id/reviews", h.GetProductReviews, middleware.ErrorHandler)
		api.GET("/products/:id", h.GetProductDetail, middleware.ErrorHandler)
		api.GET("/categories", h.GetCategories, middleware.ErrorHandler)
		api.GET("/categories/:id", h.GetCategory, middleware.ErrorHandler)
		api.GET("/sub-categories", h.GetSubCategories, middleware.ErrorHandler)
		api.GET("/sub-categories/:id", h.GetSubcategory, middleware.ErrorHandler)
		api.GET("/complementary-categories", h.GetComplementaryCategories, middleware.ErrorHandler)
		api.GET("/complementary-categories/:id", h.GetComplementaryCategory, middleware.ErrorHandler)
		api.GET("/cities", h.GetCities, middleware.ErrorHandler)
		api.GET("/couriers", h.GetAllCourier, middleware.ErrorHandler)
		api.GET("/vouchers/merchant/:id", h.GetMerchantVoucherByMerchantId, middleware.ErrorHandler)
		api.GET("/vouchers/marketplace", h.GetAllMarketplaceActiveVoucher, middleware.ErrorHandler)
		api.POST("/shipments", h.GetShippingCost, middleware.ErrorHandler)
		api.GET("/seller-categories/:username", h.GetSellerPageCategories, middleware.ErrorHandler)

		// need auth
		api.GET("/carts/hover", middleware.AuthorizeJWT, h.GetCartHoverItems, middleware.ErrorHandler)
		api.GET("/carts", middleware.AuthorizeJWT, h.GetCartItems, middleware.ErrorHandler)
		api.POST("/carts", middleware.AuthorizeJWT, h.CreateCartItem, middleware.ErrorHandler)
		api.PATCH("/carts/:id", middleware.AuthorizeJWT, h.UpdateCartItem, middleware.ErrorHandler)
		api.DELETE("/carts/:id", middleware.AuthorizeJWT, h.DeleteCartItem, middleware.ErrorHandler)
		api.PATCH("/users", middleware.AuthorizeJWT, h.EditUserDetail, middleware.ErrorHandler)
		api.GET("/wallets", middleware.AuthorizeJWT, h.GetWalletDetail, middleware.ErrorHandler)
		api.GET("/wallets/histories", middleware.AuthorizeJWT, h.GetWalletTransactionHistories, middleware.ErrorHandler)
		api.GET("/wallets/transactions", middleware.AuthorizeJWT, h.GetWalletTransactionDetail, middleware.ErrorHandler)
		api.POST("/wallets/topup", h.ValidateTopUpWithSLP, middleware.ErrorHandler)
		api.POST("/users/auth/check", middleware.AuthorizeJWT, h.CheckPasswordAuth, middleware.ErrorHandler)
		api.POST("/merchants", middleware.AuthorizeJWT, h.RegisterMerchant, middleware.ErrorHandler)
		api.GET("/merchants", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetMerchantProfile, middleware.ErrorHandler)
		api.PUT("/merchants/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.UpdateMerchantProfile, middleware.ErrorHandler)
		api.PATCH("/users/wallet", middleware.AuthorizeJWT, h.ChangeWalletPin, middleware.ErrorHandler)
		api.GET("/users/wallet/check", middleware.AuthorizeJWT, h.IsUserHasWallet, middleware.ErrorHandler)
		api.POST("/users/wallet", middleware.AuthorizeJWT, h.CreateUserWallet, middleware.ErrorHandler)
		api.POST("/wishlists", middleware.AuthorizeJWT, h.Favorite, middleware.ErrorHandler)
		api.GET("/products/:id/user-favorite", middleware.AuthorizeJWT, h.CheckFavorite, middleware.ErrorHandler)
		api.POST("/otp", middleware.AuthorizeJWT, h.GenerateOTP, middleware.ErrorHandler)
		api.GET("/users/sealabspay", middleware.AuthorizeJWT, h.GetUserSealabsAccount, middleware.ErrorHandler)
		api.POST("/users/sealabspay", middleware.AuthorizeJWT, h.AddSealabsPayAccount, middleware.ErrorHandler)
		api.PATCH("/users/sealabspay/:id", middleware.AuthorizeJWT, h.CheckSealabsOwner, h.SetDefaultSealabsAccount, middleware.ErrorHandler)
		api.DELETE("/users/sealabspay/:id", middleware.AuthorizeJWT, h.CheckSealabsOwner, h.DeleteSealabsAccount, middleware.ErrorHandler)
		api.GET("/sellers/vouchers/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.ValidateCurrentSellerVoucher, h.GetMerchantVoucherDetail, middleware.ErrorHandler)
		api.GET("/sellers/vouchers", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetMerchantVouchers, middleware.ErrorHandler)
		api.POST("/sellers/vouchers", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.CreateMerchantVoucher, middleware.ErrorHandler)
		api.POST("/users/sealabspay/validate", h.ValidateCreateSlpAccount, middleware.ErrorHandler)
		api.POST("/users/sealabspay/topup", middleware.AuthorizeJWT, h.TopUpWithSLP, middleware.ErrorHandler)
		api.PATCH("/users/sealabspay", middleware.AuthorizeJWT, h.SetDefaultSealabsAccount, middleware.ErrorHandler)
		api.GET("/sellers/products", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetProductsSellerAdmin, middleware.ErrorHandler)
		api.DELETE("/sellers/vouchers/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.ValidateCurrentSellerVoucher, h.DeleteMerchantVoucher, middleware.ErrorHandler)
		api.PATCH("/sellers/vouchers/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.ValidateCurrentSellerVoucher, h.UpdateMerchantVoucher, middleware.ErrorHandler)
		api.GET("/sellers/promos", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetProductPromos, middleware.ErrorHandler)
		api.GET("/sellers/promos/available-products", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetAvailableMerchantProductForPromo, middleware.ErrorHandler)
		api.POST("/sellers/promos", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.CreateProductPromotion, middleware.ErrorHandler)
		api.GET("/sellers/promos/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetProductPromoDetail, middleware.ErrorHandler)
		api.PUT("/sellers/promos/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.EditProductPromo, middleware.ErrorHandler)
		api.DELETE("/sellers/promos/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.EndMerchantPromo, middleware.ErrorHandler)
		api.DELETE("/sellers/promos/:id/promo-items", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.DeletePromoItems, middleware.ErrorHandler)
		api.POST("/sellers/promos/:id/promo-items", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.InsertPromoItems, middleware.ErrorHandler)
		api.POST("/otp/auth", middleware.AuthorizeJWT, h.AuthOTP, middleware.ErrorHandler)
		api.POST("/step-up/pin", middleware.AuthorizeJWT, h.ValidatePin, middleware.ErrorHandler)
		api.GET("/address", middleware.AuthorizeJWT, h.GetAddress, middleware.ErrorHandler)
		api.POST("/address", middleware.AuthorizeJWT, h.AddNewAddress, middleware.ErrorHandler)
		api.PATCH("/address", middleware.AuthorizeJWT, h.EditAddress, middleware.ErrorHandler)
		api.DELETE("/address/:id", middleware.AuthorizeJWT, h.DeleteAddress, middleware.ErrorHandler)
		api.GET("/address/default", middleware.AuthorizeJWT, h.GetDefaultAddress, middleware.ErrorHandler)
		api.PATCH("/address/default/:id", middleware.AuthorizeJWT, h.SetDefaultAddress, middleware.ErrorHandler)
		api.POST("/products", middleware.AuthorizeJWT, h.AddNewProduct, middleware.ErrorHandler)
		api.PUT("/products/:id", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.EditProduct, middleware.ErrorHandler)
		api.DELETE("/products/:id", middleware.AuthorizeJWT, h.DeleteProduct, middleware.ErrorHandler)
		api.POST("/orders/:id/reviews", middleware.AuthorizeJWT, h.CreateReview, middleware.ErrorHandler)
		api.POST("/checkout/items", middleware.AuthorizeJWT, h.AddCheckOutItems, middleware.ErrorHandler)
		api.GET("/checkout/items", middleware.AuthorizeJWT, h.GetCheckOutItems, middleware.ErrorHandler)
		api.GET("/sellers/categories", middleware.AuthorizeJWT, h.GetSellerCategories, middleware.ErrorHandler)
		api.GET("/sellers/categories/:id", middleware.AuthorizeJWT, h.GetSellerCategory, middleware.ErrorHandler)
		api.POST("/sellers/categories", middleware.AuthorizeJWT, h.CreateSellerCategory, middleware.ErrorHandler)
		api.PUT("/sellers/categories/:id", middleware.AuthorizeJWT, h.UpdateSellerCategory, middleware.ErrorHandler)
		api.DELETE("/sellers/categories/:id", middleware.AuthorizeJWT, h.DeleteSellerCategory, middleware.ErrorHandler)
		api.GET("/sellers/categories/available-products", middleware.AuthorizeJWT, h.GetAvailableProducts, middleware.ErrorHandler)
		api.POST("/sellers/categories/:id/items", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.AddProductsToSellerCategory, middleware.ErrorHandler)
		api.DELETE("/sellers/categories/items", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.RemoveProductFromSellerCategory, middleware.ErrorHandler)
		api.GET("/sellers/orders", middleware.AuthorizeJWT, middleware.AuthorizeSeller, h.GetMerchantOrders, middleware.ErrorHandler)
		api.POST("/checkout/pay/wallet", middleware.AuthorizeStepUpJWT, h.PayCheckoutWithWallet, middleware.ErrorHandler)
		api.POST("/checkout/proceed", middleware.AuthorizeJWT, h.ProceedOrder, middleware.ErrorHandler)
		api.GET("/paychecks", middleware.AuthorizeJWT, h.GetMerchantPaycheck, middleware.ErrorHandler)
		api.POST("/paychecks/withdraw", middleware.AuthorizeJWT, h.WithDrawMerchantIncome, middleware.ErrorHandler)
		api.GET("/paychecks/history", middleware.AuthorizeJWT, h.GetWithdrawHistories, middleware.ErrorHandler)
		api.PUT("/order/:id", middleware.AuthorizeJWT, h.UpdateOrderStatus, middleware.ErrorHandler)

		// system admin
		api.POST("/categories", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.CreateCategory, middleware.ErrorHandler)
		api.POST("/sub-categories", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.CreateSubCategory, middleware.ErrorHandler)
		api.POST("/complementary-categories", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.CreateComplementaryCategory, middleware.ErrorHandler)
		api.PATCH("/categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.UpdateCategory, middleware.ErrorHandler)
		api.PATCH("/sub-categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.UpdateSubCategory, middleware.ErrorHandler)
		api.PATCH("/complementary-categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.UpdateComplementaryCategory, middleware.ErrorHandler)
		api.DELETE("/categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.DeleteCategory, middleware.ErrorHandler)
		api.DELETE("/sub-categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.DeleteSubCategory, middleware.ErrorHandler)
		api.DELETE("/complementary-categories/:id", middleware.AuthorizeJWT, middleware.AuthorizeAdmin, h.DeleteComplementaryCategory, middleware.ErrorHandler)
	}

	return router
}
