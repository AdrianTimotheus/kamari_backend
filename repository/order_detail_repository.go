package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type OrderDetailsRepository interface {
	GetOrderDetails(productVariantId uint64) ([]model.OrderDetail, error)
}

type orderDetailsRepository struct {
	db *gorm.DB
}

type ODRConfig struct {
	DB *gorm.DB
}

func NewOrderDetailRepository(c ODRConfig) OrderDetailsRepository {
	return orderDetailsRepository{
		db: c.DB,
	}
}

func (o orderDetailsRepository) GetOrderDetails(productVariantId uint64) ([]model.OrderDetail, error) {
	var order []model.OrderDetail
	err := o.db.
		Preload("ProductReview").
		Where("product_variant_id = ?", productVariantId).
		Find(&order).Error

	return order, err
}
