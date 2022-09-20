package server

import (
	"kamari.com/db"
	"kamari.com/repository"
	"kamari.com/service"
	"log"
)

func Init() {

	provinceRepository := repository.NewProvinceRepository(repository.PRConfig{DB: db.Get()})
	provinceService := service.NewProvinceService(service.PSConfig{ProvinceRepository: provinceRepository})

	r := NewRouter(&RouterConfig{ProvinceService: provinceService})

	log.Fatal(r.Run())
}
