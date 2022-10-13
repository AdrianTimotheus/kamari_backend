package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type CityRepository interface {
	GetCity(id uint64) (model.City, error)
	GetCityPerProvince() ([]model.City, error)
}

type cityRepository struct {
	db *gorm.DB
}

type CiRConfig struct {
	DB *gorm.DB
}

func NewCityRepository(c CRConfig) CityRepository {
	return cityRepository{
		db: c.DB,
	}
}

func (c cityRepository) GetCity(id uint64) (model.City, error) {
	var city model.City
	err := c.db.Joins("Province").Where("cities.id = ?", id).First(&city).Error
	if err != nil {
		return city, err
	}

	return city, nil
}

func (c cityRepository) GetCityPerProvince() ([]model.City, error) {
	var models []model.City
	err := c.db.Joins("Province").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
