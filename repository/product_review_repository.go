package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProductReviewRepository interface {
	GetProductReview(productId uint64) ([]model.ProductReview, error)
}

type productReviewRepository struct {
	db *gorm.DB
}

type PRRConfig struct {
	DB *gorm.DB
}

func NewProductReviewRepository(c PRRConfig) ProductReviewRepository {
	return productReviewRepository{
		db: c.DB,
	}
}

func (p productReviewRepository) GetProductReview(productId uint64) ([]model.ProductReview, error) {
	var (
		reviews          []model.ProductReview
		prodId           []int
		productVariantId []int
	)

	err := p.db.
		Table("product_variants").
		Select("id").
		Where("product_id = ?", productId).
		Scan(&prodId).Error
	if err != nil {
		return nil, err
	}

	err = p.db.
		Table("order_details").
		Select("id").
		Where("product_variant_id IN ?", prodId).
		Scan(&productVariantId).Error
	if err != nil {
		return nil, err
	}

	err = p.db.
		Preload("ReviewPhotos").
		Where("order_detail_id IN ?", productVariantId).
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
