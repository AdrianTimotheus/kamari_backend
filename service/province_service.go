package service

import (
	"kamari.com/config"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type ProvinceService interface {
	GetAllProvinces() ([]model.Province, error)
}

type provinceService struct {
	provinceRepository repository.ProvinceRepository
	appConfig          config.AppConfig
}

type PSConfig struct {
	ProvinceRepository repository.ProvinceRepository
	AppConfig          config.AppConfig
}

func NewProvinceService(p PSConfig) ProvinceService {
	return provinceService{
		provinceRepository: p.ProvinceRepository,
		appConfig:          p.AppConfig,
	}
}

func (p provinceService) GetAllProvinces() ([]model.Province, error) {
	result, err := p.provinceRepository.GetAllProvinces()
	if err != nil {
		return nil, httperror.InternalServerError(err.Error())
	}
	return result, nil
}
