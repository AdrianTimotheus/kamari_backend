package server

import (
	"kamari.com/db"
	"kamari.com/repository"
	"kamari.com/service"
	"log"
)

func Init() {
	//REPO
	userRepository := repository.NewUserRepository(repository.URConfig{
		DB: db.Get(),
	})

	merchantRepository := repository.NewMerchantRepository(repository.MRConfig{
		DB: db.Get(),
	})

	productRepository := repository.NewProductRepository(repository.ProdConfig{
		DB: db.Get(),
	})

	merchantFollowerRepository := repository.NewMFRepository(repository.MFConfig{
		DB: db.Get(),
	})

	cityRepository := repository.NewCityRepository(repository.CRConfig{
		DB: db.Get(),
	})

	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{
		DB: db.Get(),
	})

	marketplaceRepository := repository.NewMarketplaceRepository(repository.MpRConfig{
		DB: db.Get(),
	})

	cartRepository := repository.NewCartRepository(repository.CRConfig{
		DB: db.Get(),
	})

	//SERVICE
	userService := service.NewUserService(service.USConfig{
		UserRepository: userRepository,
	})

	merchantService := service.NewMerchantService(service.MSConfig{
		MerchantRepository:         merchantRepository,
		CityRepository:             cityRepository,
		ProductRepository:          productRepository,
		MerchantFollowerRepository: merchantFollowerRepository,
	})

	provinceService := service.NewProvinceService(service.PSConfig{
		ProvinceRepository: provinceRepository,
	})

	marketplaceService := service.NewMarketplaceService(service.MpSConfig{
		MarketplaceRepository: marketplaceRepository,
	})

	cartService := service.NewCartService(service.CSConfig{
		CartRepository: cartRepository,
	})

	productService := service.NewProductService(service.PConfig{
		ProductRepository:  productRepository,
		MerchantRepository: merchantRepository,
		CityRepository:     cityRepository,
	})

	r := NewRouter(RouterConfig{
		UserService:        userService,
		ProvinceService:    provinceService,
		CartService:        cartService,
		MerchantService:    merchantService,
		MarketplaceService: marketplaceService,
		ProductService:     productService,
	})

	log.Fatal(r.Run())
}
