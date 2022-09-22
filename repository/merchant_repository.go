package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type MerchantRepository interface {
	GetMerchantDetail(username string) (model.Merchant, error)
}

type merchantRepository struct {
	db *gorm.DB
}

type MRConfig struct {
	DB *gorm.DB
}

func NewMerchantRepository(c MRConfig) MerchantRepository {
	return merchantRepository{
		db: c.DB,
	}
}

func (m merchantRepository) GetMerchantDetail(username string) (model.Merchant, error) {
	var merchant model.Merchant
	err := m.db.
		Joins("User").
		Joins("UserAddress").
		Where(`"User"."username" = ?`, username).
		First(&merchant).Error

	if err != nil {
		return merchant, err
	}

	return merchant, nil
}
