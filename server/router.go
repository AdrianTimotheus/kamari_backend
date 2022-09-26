package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kamari.com/handler"
	"kamari.com/middleware"
	"kamari.com/service"
	"net/http"
	"os"
	"time"
)

type RouterConfig struct {
	UserService     service.UserService
	ProvinceService service.ProvinceService
	CartService     service.CartService
	MerchantService service.MerchantService
	MarketplaceService service.MarketplaceService
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
		UserService:     c.UserService,
		ProvinceService: c.ProvinceService,
		CartService:     c.CartService,
		MerchantService: c.MerchantService,
		MarketplaceService: c.MarketplaceService,
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("route not found, go to %v/docs for documentation\n", os.Getenv("BASE_URL")),
		})
	})

	api := router.Group("/api")
	{
		api.Static("/docs", "swaggerui")
		api.POST("/users", h.Register, middleware.ErrorHandler)
		api.POST("/login", h.Login, middleware.ErrorHandler)
		api.GET("/provinces", h.GetAllProvinces, middleware.ErrorHandler)
		api.GET("/carts/hover", middleware.AuthorizeJWT, h.GetCartHoverItems, middleware.ErrorHandler)
		api.GET("/sellers/:username", h.GetSellerDetail, middleware.ErrorHandler)
		api.GET("/merchant/:username/products", h.GetSellerProduct, middleware.ErrorHandler)
		api.GET("/marketplace/promo/banners", h.GetMarketplacePromoBanner, middleware.ErrorHandler)
	}

	return router
}
