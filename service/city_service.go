package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type CityService interface {
	GetListCityPerProvince() (map[string][]model.City, error)
}

type cityService struct {
	cityRepository repository.CityRepository
}

type CiSConfig struct {
	CityRepository repository.CityRepository
}

func NewCityService(c CiSConfig) CityService {
	return cityService{cityRepository: c.CityRepository}
}

func (c cityService) GetListCityPerProvince() (map[string][]model.City, error) {
	models, err := c.cityRepository.GetCityPerProvince()

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	temp := make(map[string][]model.City)

	for _, data := range models {
		tempData, ok := temp[data.Province.Name]
		if !ok {
			tempData = append(tempData, data)
			temp[data.Province.Name] = tempData
		} else {
			temp[data.Province.Name] = append(tempData, data)
		}
	}
	return temp, nil
}
