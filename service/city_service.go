package service

import (
	"kamari.com/repository"
)

type CityService interface {
}

type cityService struct {
	merchantRepository repository.MerchantRepository
}

type CSConfig struct {
	MerchantRepository repository.MerchantRepository
}

func NewCityService(c MSConfig) CityService {
	return cityService{
		merchantRepository: c.MerchantRepository,
	}
}
