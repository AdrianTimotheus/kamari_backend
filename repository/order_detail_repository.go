package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type OrderDetailsRepository interface {
	GetOrderDetails(id uint64) (*model.OrderDetail, error)
	GetOrdererUser(orderDetailId uint64) (*model.User, error)
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

func (o orderDetailsRepository) GetOrdererUser(orderDetailId uint64) (*model.User, error) {
	var orderDetail model.OrderDetail
	err := o.db.
		Preload("Order").
		Preload("Order.Transaction").
		Preload("Order.Transaction.User").
		Where("id = ?", orderDetailId).
		First(&orderDetail).Error
	if err != nil {
		return nil, err
	}

	return &orderDetail.Order.Transaction.User, nil
}
