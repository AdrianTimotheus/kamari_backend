package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type MerchantFollowerRepository interface {
	GetMerchantFollowerCount(merchantId uint64) (uint64, error)
}

type merchantFollowerRepository struct {
	db *gorm.DB
}

type MFConfig struct {
	DB *gorm.DB
}

func NewMFRepository(c MFConfig) MerchantFollowerRepository {
	return merchantFollowerRepository{
		db: c.DB,
	}
}

func (m merchantFollowerRepository) GetMerchantFollowerCount(merchantId uint64) (uint64, error) {
	var count int64
	err := m.db.
		Model(&model.MerchantFollower{}).
		Where("merchant_id = ?", merchantId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return uint64(count), nil
}
