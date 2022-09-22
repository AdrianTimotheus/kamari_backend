package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProductRepository interface {
	GetMerchantProduct(merchantId uint64) ([]model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

type ProdConfig struct {
	DB *gorm.DB
}

func NewProductRepository(c ProdConfig) ProductRepository {
	return productRepository{
		db: c.DB,
	}
}

func (c productRepository) GetMerchantProduct(merchantId uint64) ([]model.Product, error) {
	var product []model.Product
	err := c.db.
		Where("merchant_id = ?", merchantId).
		Find(&product).Error

	if err != nil {
		return product, err
	}
	return product, nil
}
