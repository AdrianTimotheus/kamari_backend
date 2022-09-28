package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type OrderDetailsRepository interface {
	GetOrderDetails(id uint64) (*model.OrderDetail, error)
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

func (o orderDetailsRepository) GetOrderDetails(id uint64) (*model.OrderDetail, error) {
	var order *model.OrderDetail
	err := o.db.
		Preload("ProductReview").
		Where("id = ?", id).
		Find(&order).Error

	return order, err
}
