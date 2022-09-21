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
	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{DB: db.Get()})

	userService := service.NewUserService(service.USConfig{
		UserRepository: userRepository,
	})
	provinceService := service.NewProvinceService(service.PSConfig{ProvinceRepository: provinceRepository})

	r := NewRouter(RouterConfig{
		UserService:     userService,
		ProvinceService: provinceService,
	})

	log.Fatal(r.Run())
}
