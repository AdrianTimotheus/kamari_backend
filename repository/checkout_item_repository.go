package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type CheckoutItemRepository interface {
	CreateCheckoutRecords(models []model.CheckoutItem) error
	GetCheckoutItemByUser(userid uint64) ([]model.CheckoutItem, error)
	GetSupportedCourierForMerchant(merchantId uint64, listOfProductId []uint64) ([]model.Courier, error)
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

func (c checkoutItemRepository) GetSupportedCourierForMerchant(merchantId uint64, listOfProductId []uint64) ([]model.Courier, error) {
	//select product_couriers.courier_id from product_couriers where
	//	product_couriers.product_id in (7,10) group by product_couriers.courier_id  having count(product_couriers.courier_id)=2
	var output []model.Courier
	var pCModel model.ProductCourier
	var merchantCourierModel model.MerchantCourier

	var merchantSupportedCourierIds []uint64
	c.db.Model(&merchantCourierModel).
		Select("courier_id").Where("merchant_id=?", merchantId).Scan(&merchantSupportedCourierIds)

	var productSupportedCourierIds []uint64
	c.db.Model(&pCModel).
		Select("courier_id").
		Where("product_id IN ?", listOfProductId).
		Group("courier_id").
		Having("count(courier_id)=?", len(listOfProductId)).
		Where("courier_id IN?", merchantSupportedCourierIds).Find(&productSupportedCourierIds)

	err := c.db.Model(&output).Where("id IN ?", productSupportedCourierIds).Find(&output).Error
	return output, err
}
