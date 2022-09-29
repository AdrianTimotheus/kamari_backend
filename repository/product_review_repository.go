package repository

import (
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/model"
)

type ProductReviewRepository interface {
	GetProductReview(productId uint64, filters *helper.Filters) ([]model.ProductReview, *helper.Metadata, error)
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

func (p productReviewRepository) GetProductReview(productId uint64, filters *helper.Filters) ([]model.ProductReview, *helper.Metadata, error) {
	var (
		reviews          []model.ProductReview
		prodId           []int
		productVariantId []int
		count            int64
	)

	err := p.db.
		Table("product_variants").
		Select("id").
		Where("product_id = ?", productId).
		Scan(&prodId).Error
	if err != nil {
		return nil, nil, err
	}

	err = p.db.
		Table("order_details").
		Select("id").
		Where("product_variant_id IN ?", prodId).
		Scan(&productVariantId).Error
	if err != nil {
		return nil, nil, err
	}

	err = p.db.
		Preload("ReviewPhotos").
		Where("order_detail_id IN ?", productVariantId).
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&reviews).Error
	if err != nil {
		return nil, nil, err
	}

	err = p.db.
		Table("product_reviews").
		Select("COUNT(*) OVER() AS COUNT").
		Limit(1).
		Scan(&count).Error
	if err != nil {
		return nil, nil, err
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return reviews, &metadata, nil
}
