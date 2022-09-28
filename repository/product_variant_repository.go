package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProductVariantRepository interface {
	GetProductVariants(productId uint64) ([]model.ProductVariant, error)
	GetProductVariantById(variantId uint64) ([]model.ProductVariant, error)
}

type productVariantRepository struct {
	db *gorm.DB
}

type PVConfig struct {
	DB *gorm.DB
}

func NewProductVariantRepository(c PVConfig) ProductVariantRepository {
	return productVariantRepository{
		db: c.DB,
	}
}

func (p productVariantRepository) GetProductVariants(productId uint64) ([]model.ProductVariant, error) {
	var variants []model.ProductVariant

	err := p.db.Where("product_id = ?", productId).Find(&variants).Error
	if err != nil {
		return nil, err
	}

	return variants, nil
}

func (p productVariantRepository) GetProductVariantById(variantId uint64) ([]model.ProductVariant, error) {
	var variants []model.ProductVariant

	err := p.db.Where("id = ?", variantId).Find(&variants).Error
	if err != nil {
		return nil, err
	}

	return variants, nil
}
