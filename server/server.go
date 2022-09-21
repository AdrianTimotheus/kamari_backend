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

	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{
		DB: db.Get(),
	})

	marketplaceRepository := repository.NewMarketplaceRepository(repository.MpRConfig{
		DB: db.Get(),
	})

	//SERVICE
	userService := service.NewUserService(service.USConfig{
		UserRepository: userRepository,
	})

	provinceService := service.NewProvinceService(service.PSConfig{
		ProvinceRepository: provinceRepository,
	})

	marketplaceService := service.NewMarketplaceService(service.MpSConfig{
		MarketplaceRepository: marketplaceRepository,
	})

	r := NewRouter(RouterConfig{
		UserService:        userService,
		ProvinceService:    provinceService,
		MarketplaceService: marketplaceService,
	})

	log.Fatal(r.Run())
}
