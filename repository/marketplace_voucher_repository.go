package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
	"time"
)

type MarketplaceVoucherRepository interface {
	GetAllActiveVoucher() ([]model.MarketplaceVoucher, error)
}

type marketplaceVoucherRepository struct {
	db *gorm.DB
}

type MVCRConfig struct {
	DB *gorm.DB
}

func NewMarketplaceVoucherRepository(c MVCRConfig) MarketplaceVoucherRepository {
	return &marketplaceVoucherRepository{db: c.DB}
}

func (m marketplaceVoucherRepository) GetAllActiveVoucher() ([]model.MarketplaceVoucher, error) {
	now := time.Now().Format("2006-01-02")
	var models []model.MarketplaceVoucher
	err := m.db.
		//Joins("PromoType").
		Where("start_date >= ?", now).
		Where("end_date <= ?", now).
		Where("stock_left > 0").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
