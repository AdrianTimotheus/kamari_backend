package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProvinceRepository interface {
	GetAllProvinces() ([]model.Province, error)
}

type provinceRepository struct {
	db *gorm.DB
}

type PRConfig struct {
	DB *gorm.DB
}

func NewProvinceRepository(c PRConfig) ProvinceRepository {
	return provinceRepository{db: c.DB}
}

func (p provinceRepository) GetAllProvinces() ([]model.Province, error) {
	var provinces []model.Province

	result := p.db.Find(&provinces)
	if result.Error != nil {
		return nil, result.Error
	}

	return provinces, nil
}
