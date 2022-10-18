package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type CheckoutItemRepository interface {
	CreateCheckoutRecords(models []model.CheckoutItem) error
}

type checkoutItemRepository struct {
	db *gorm.DB
}

type CIRConfig struct {
	DB *gorm.DB
}

func NewCheckoutRepository(c CIRConfig) CheckoutItemRepository {
	return &checkoutItemRepository{db: c.DB}
}

func (c checkoutItemRepository) CreateCheckoutRecords(models []model.CheckoutItem) error {
	c.db.Exec("DELETE FROM checkout_items WHERE user_id=?", models[0].UserId)
	return c.db.Create(models).Error
}
