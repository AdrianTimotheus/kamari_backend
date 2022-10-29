package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/model"
)

type OrderRepository interface {
	GetMerchantOrders(merchantId uint64, keyword string, filters *helper.Filters) ([]model.SellerOrder, *helper.Metadata, error)
	GetMerchantOrderSummary(merchantId uint64) ([]model.SellerOrder, error)
	CreateOrder(user model.User, order []model.Order, item []dto.ItemToProceed, listOfProductVariant []uint64) error
	GetOrderItemToProcess(userId uint64) ([]model.CheckoutItem, error)
	CountUserCheckoutItem(listOfProductVariant []uint64, userId uint64) int64
	GetCurrentOrder(orderId uint64) (*model.Order, error)
	UpdateCurrentOrderId(orderModel *model.Order) error
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

func (o orderRepository) GetMerchantOrders(merchantId uint64, keyword string, filters *helper.Filters) ([]model.SellerOrder, *helper.Metadata, error) {
	var (
		orders []model.SellerOrder
		count  int64
	)

	invoiceSearch := fmt.Sprintf("%%%s%%", keyword)
	res := o.db.
		Table("orders").
		Select(
			`"orders"."id" AS id`,
			`"orders"."invoice" AS invoice`,
			`"orders"."user_total_payment" AS user_total_payment`,
			`"orders"."transaction_id" AS transaction_id`,
			`"orders"."merchant_id" AS merchant_id`,
			`"c"."name" AS courier_name`,
			`"orders"."withdraw_status" AS withdraw_status`,
			`"os"."status" AS status`,
			`"u"."username" AS user_username`,
			`"u"."pic_url" AS user_pic_url`,
			`"u"."phone" AS user_phone`,
			`"u"."email" AS user_email`,
			`"u"."name" AS user_name`,
		).
		Joins(`LEFT JOIN order_statuses os on orders.order_status_id = os.id`).
		Joins(`LEFT JOIN users u on orders.user_id = u.id`).
		Joins(`LEFT JOIN couriers c on orders.courier_id = c.id`).
		Where(`"orders"."merchant_id" = ?`, merchantId).
		Where(`"orders"."invoice" ILIKE ?`, invoiceSearch)
	if filters.OrderStatus != nil {
		res = res.Where(`"os"."status" = ?`, *filters.OrderStatus)
	}

	if err := res.Find(&orders).Count(&count).Error; err != nil {
		return nil, nil, err
	}

	res = res.
		Order("id DESC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&orders)
	if res.Error != nil {
		return nil, nil, nil
	}

	for i, order := range orders {
		od, err := o.getOrderDetails(order.Id)
		if err != nil {
			return nil, nil, err
		}

		orders[i].OrderDetail = od
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return orders, &metadata, nil
}

func (o orderRepository) GetMerchantOrderSummary(merchantId uint64) ([]model.SellerOrder, error) {
	var orders []model.SellerOrder

	res := o.db.
		Table("orders").
		Select(
			`"orders"."id" AS id`,
			`"orders"."created_at" AS order_created_at`,
			`"orders"."invoice" AS invoice`,
			`"orders"."user_total_payment" AS user_total_payment`,
			`"orders"."transaction_id" AS transaction_id`,
			`"orders"."merchant_id" AS merchant_id`,
			`"orders"."withdraw_status" AS withdraw_status`,
			`"os"."status" AS status`,
		).
		Joins(`LEFT JOIN order_statuses os on orders.order_status_id = os.id`).
		Where(`"orders"."merchant_id" = ?`, merchantId).
		Find(&orders)

	if res.Error != nil {
		return nil, res.Error
	}
	return orders, nil
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

func (o orderRepository) getOrderDetails(orderId uint64) ([]model.SellerOrderDetail, error) {
	var orderDetails []model.SellerOrderDetail
	res := o.db.
		Table("order_details").
		Select(
			`"order_details"."id" AS id`,
			`"order_details"."quantity" AS quantity`,
			`"p"."id" AS product_id`,
			`"p"."name" AS product_name`,
			`"p"."thumbnail_url" AS thumbnail_url`,
			`"pv"."main_variant" AS main_variant`,
			`"pv"."main_variant_item" AS main_variant_item`,
			`"pv"."sub_variant" AS sub_variant`,
			`"pv"."sub_variant_item" AS sub_variant_item`,
		).
		Joins(`LEFT JOIN product_variants pv on order_details.product_variant_id = pv.id`).
		Joins(`LEFT JOIN products p on pv.product_id = p.id`).
		Where(`"order_details"."order_id" = ?`, orderId).
		Scan(&orderDetails)
	if res.Error != nil {
		return nil, res.Error
	}

	return orderDetails, nil
}

func (o orderRepository) GetCurrentOrder(orderId uint64) (*model.Order, error) {
	orderModel := model.Order{
		Id: orderId,
	}
	err := o.db.Joins("OrderStatus").
		Joins("Merchant").
		First(&orderModel).Error
	if err != nil {
		return nil, err
	}
	return &orderModel, nil
}

func (o orderRepository) UpdateCurrentOrderId(orderModel *model.Order) error {
	fmt.Println("HEY YOOOOO:", orderModel.OrderStatusId)
	err := o.db.Table("orders").Where("id=?", orderModel.Id).Update("order_status_id", orderModel.OrderStatusId).Error
	if err != nil {
		return err
	}
	return nil
}
