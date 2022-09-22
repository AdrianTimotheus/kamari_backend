package server

import (
	"kamari.com/db"
	"kamari.com/repository"
	"kamari.com/service"
	"log"
)

func Init() {
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
	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{DB: db.Get()})

	userService := service.NewUserService(service.USConfig{
		UserRepository: userRepository,
	})
	merchantService := service.NewMerchantService(service.MSConfig{
		MerchantRepository:         merchantRepository,
		CityRepository:             cityRepository,
		ProductRepository:          productRepository,
		MerchantFollowerRepository: merchantFollowerRepository,
	})
	provinceService := service.NewProvinceService(service.PSConfig{ProvinceRepository: provinceRepository})

	r := NewRouter(RouterConfig{
		UserService:     userService,
		ProvinceService: provinceService,
		MerchantService: merchantService,
	})

	log.Fatal(r.Run())
}
