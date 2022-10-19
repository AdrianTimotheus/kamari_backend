package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type CheckoutItemRepository interface {
	CreateCheckoutRecords(models []model.CheckoutItem) error
	GetCheckoutItemByUser(userid uint64) ([]model.CheckoutItem, error)
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

func (c checkoutItemRepository) GetCheckoutItemByUser(userId uint64) ([]model.CheckoutItem, error) {
	var models []model.CheckoutItem
	err := c.db.
		Where("checkout_items.user_id=?", userId).
		Joins("ProductVariant").
		Joins("Merchant").
		Joins("User").
		Find(&models).Error

	if err != nil {
		return nil, err
	}
	return models, nil
}
