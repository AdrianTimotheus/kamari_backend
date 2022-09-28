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
	UserService              service.UserService
	ProvinceService          service.ProvinceService
	MerchantService          service.MerchantService
	CartService              service.CartService
	MarketplaceService       service.MarketplaceService
	ProductService           service.ProductService
	WalletService            service.WalletService
	CategoryService          service.CategoryService
	WalletTransactionService service.WalletTransactionService
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
		UserService:              c.UserService,
		ProvinceService:          c.ProvinceService,
		MerchantService:          c.MerchantService,
		ProductService:           c.ProductService,
		CartService:              c.CartService,
		MarketplaceService:       c.MarketplaceService,
		WalletService:            c.WalletService,
		CategoryService:          c.CategoryService,
		WalletTransactionService: c.WalletTransactionService,
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("route not found, go to %v/docs for documentation\n", os.Getenv("BASE_URL")),
		})
	})

	api := router.Group("/api")
	{
		// guest
		api.Static("/docs", "swaggerui")
		api.POST("/users", h.Register, middleware.ErrorHandler)
		api.POST("/login", h.Login, middleware.ErrorHandler)
		api.GET("/provinces", h.GetAllProvinces, middleware.ErrorHandler)
		api.GET("/products", h.GetProducts, middleware.ErrorHandler)
		api.GET("/categories", h.GetCategories, middleware.ErrorHandler)
		api.GET("/merchants/:username", h.GetMerchantDetail, middleware.ErrorHandler)
		api.GET("/products/:id", h.GetProductDetail, middleware.ErrorHandler)
		api.GET("/sellers/:username", h.GetMerchantDetail, middleware.ErrorHandler)
		api.GET("/sellers/:username/products", h.GetSellerProduct, middleware.ErrorHandler)
		api.GET("/marketplace/promo/banners", h.GetMarketplacePromoBanner, middleware.ErrorHandler)
		api.GET("/products/recommend", h.GetRecommendedProduct, middleware.ErrorHandler)
		// need auth
		api.GET("/carts/hover", middleware.AuthorizeJWT, h.GetCartHoverItems, middleware.ErrorHandler)
		api.GET("/wallets", middleware.AuthorizeJWT, h.GetWalletDetail, middleware.ErrorHandler)
		api.GET("/wallets/histories", middleware.AuthorizeJWT, h.GetWalletTransactionHistories, middleware.ErrorHandler)
		api.GET("/wallets/transactions", middleware.AuthorizeJWT, h.GetWalletTransactionDetail, middleware.ErrorHandler)
	}

	return router
}
