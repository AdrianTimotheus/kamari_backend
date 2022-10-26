package repository

import (
	"gorm.io/gorm"
)

type OrderRepository interface {
	//GetMerchantOrders(merchantId uint64, keyword string, filters *helper.Filters) ([]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

type ORConfig struct {
	DB *gorm.DB
}

func NewOrderRepository(c ORConfig) OrderRepository {
	return orderRepository{
		db: c.DB,
	}
}

//func (o orderRepository) GetMerchantOrders(merchantId uint64, keyword string, filters *helper.Filters) ([]model.Order, error) {
//	var orders []model.Order
//
//	k := fmt.Sprintf("%%%s", keyword)
//	o.db.
//		Table("orders").
//		Select(
//			`"orders"."id" AS id`,
//			"invoice",
//			"user_total_payment",
//			"merchant_income",
//			"shipping_fee",
//			"total_products_price",
//			"merchant_voucher_id",
//			"marketplace_voucher_id",
//			"dst_address_id",
//			"transaction_id",
//			"order_status_id",
//			"withdraw_status",
//			"courier_id",
//			"user_id",
//			`"orders"."created_at" AS created_at`,
//			`"orders"."updated_at" AS updated_at`,
//			`"orders"."deleted_at" AS deleted_at`,
//		).
//		Joins(`LEFT JOIN "order_details" AS "od" ON "od"."order_id" = "orders"."id"`).
//		Joins(`LEFT JOIN "product_variants" AS "pv" ON "pv"."id" = "od"."product_variant_id"`)
//}
