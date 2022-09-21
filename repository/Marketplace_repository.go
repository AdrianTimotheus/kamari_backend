package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type MarketplaceRepository interface {
	GetPromoBanner(date string) ([]model.MarketplaceVoucher, error)
}

type marketplaceRepository struct {
	db *gorm.DB
}

type MpRConfig struct {
	DB *gorm.DB
}

func NewMarketplaceRepository(config MpRConfig) MarketplaceRepository {
	return marketplaceRepository{
		db: config.DB,
	}
}

func (m marketplaceRepository) GetPromoBanner(date string) ([]model.MarketplaceVoucher, error) {
	var marketplaceVoucherBanners []model.MarketplaceVoucher

	err := m.db.
		Where("start_date <=?", date).
		Where("end_date >=?", date).
		Order("start_date desc").
		Find(&marketplaceVoucherBanners).Error

	if err != nil {
		return nil, err
	}
	return marketplaceVoucherBanners, nil
}
