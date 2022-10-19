package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type CourierRepository interface {
	GetAllCourier() ([]model.Courier, error)
}

type courierRepository struct {
	db *gorm.DB
}

type CConfig struct {
	DB *gorm.DB
}

func NewCourierRepository(c CConfig) CourierRepository {
	return &courierRepository{db: c.DB}
}

func (c courierRepository) GetAllCourier() ([]model.Courier, error) {
	var models []model.Courier
	err := c.db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
