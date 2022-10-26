package repository

import (
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/model"
)

type OrderRepository interface {
	CreateOrder(user model.User, order []model.Order, item []dto.ItemToProceed, listOfProductVariant []uint64) error
	GetOrderItemToProcess(userId uint64) ([]model.CheckoutItem, error)
	CountUserCheckoutItem(listOfProductVariant []uint64, userId uint64) int64
}

type orderRepository struct {
	db *gorm.DB
}

type ORConfig struct {
	DB *gorm.DB
}

func NewOrderRepository(c ORConfig) OrderRepository {
	return &orderRepository{
		db: c.DB,
	}
}

func (o orderRepository) CreateOrder(user model.User, order []model.Order, item []dto.ItemToProceed, listOfProductVariant []uint64) error {
	err := o.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&order).Error
		if err != nil {
			return err
		}

		var orderDetailModels []model.OrderDetail
		for _, data := range item {
			temp := model.OrderDetail{
				Quantity:         uint64(data.Quantity),
				ProductVariantId: data.VariantID,
				RealPrice:        data.RealPrice,
				DiscountedPrice:  data.Price,
			}
			for _, createdOrder := range order {
				if createdOrder.MerchantId == data.MerchantId {
					temp.OrderId = createdOrder.Id
					break
				}
			}
			orderDetailModels = append(orderDetailModels, temp)
		}
		err = tx.Create(&orderDetailModels).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM checkout_items where user_id=?", user.Id).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM carts where user_id=? AND variant_id IN ?", user.Id, listOfProductVariant).Error
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (o orderRepository) GetOrderItemToProcess(userId uint64) ([]model.CheckoutItem, error) {
	var checkoutItem []model.CheckoutItem
	err := o.db.Model(&checkoutItem).Where("user_id", userId).Find(&checkoutItem).Error
	if err != nil {
		return nil, err
	}
	return checkoutItem, nil
}

func (o orderRepository) CountUserCheckoutItem(listOfProductVariant []uint64, userId uint64) int64 {
	var modelToCheck model.CheckoutItem
	var numberOfRecords int64
	o.db.Model(&modelToCheck).
		Where("product_variant_id in ?", listOfProductVariant).
		Where("user_id=?", userId).
		Count(&numberOfRecords)
	return numberOfRecords
}
