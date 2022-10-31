package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	CreateOrder(user model.User, payload dto.CreateOrderRequest) (*dto.StatusMessageResponse, error)
	GetMerchantOrder(merchantId uint64, keyword string, filters *helper.Filters) ([]model.SellerOrder, *helper.Metadata, error)
	GetUserOrders(userId uint64, keyword string, filters *helper.Filters) ([]dto.UserOrderRes, *helper.Metadata, error)
	CreateOrderSLP(sign string) (*dto.StatusMessageResponse, error)
	GetOrderDetails(orderId, merchantId uint64) (*dto.OrderHighlightResponse, error)
	GetMerchantOrderSummary(merchantId uint64) ([]model.SellerOrder, error)
	UpwardOrderStatus(OrderId uint64, user model.User) (*dto.StatusMessageResponse, error)
	CompleteOrder(orderId uint64, user *model.User) error
}

type orderService struct {
	orderRepository repository.OrderRepository
}

type OSConfig struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(c OSConfig) OrderService {
	return &orderService{orderRepository: c.OrderRepository}
}

func (o orderService) CreateOrder(user model.User, payload dto.CreateOrderRequest) (*dto.StatusMessageResponse, error) {
	//verify product
	var listOfProductVariant []uint64
	for _, data := range payload.Products {
		listOfProductVariant = append(listOfProductVariant, data.VariantID)
	}
	numOfRecords := o.orderRepository.CountUserCheckoutItem(listOfProductVariant, user.Id)
	if numOfRecords != int64(len(listOfProductVariant)) {
		return nil, httperror.BadRequestError("Wrong Product Variant", "WRONG_PRODUCT_VARIANT")
	}

	var orderModels []model.Order
	for i, data := range payload.Orders {
		now := time.Now().Format("2006-01-02-150405")
		invoice := "INV/?/TRC-?/MRCH-?/USR-?/?"
		invoice = strings.Replace(invoice, "?", now, 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(payload.TransactionId)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(data.MerchantId)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(user.Id)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(i), 1)

		temp := model.Order{
			Invoice:                  invoice,
			UserTotalPayment:         data.UserTotalPayment,
			MerchantIncome:           data.MerchantIncome,
			MerchantVoucherId:        data.MerchantVoucherId,
			MarketplaceVoucherId:     data.MarketplaceVoucherId,
			DstAddressId:             data.DstAddressId,
			TransactionId:            payload.TransactionId,
			OrderStatusId:            1,
			WithdrawStatus:           false,
			MerchantId:               data.MerchantId,
			CourierId:                data.CourierId,
			UserId:                   user.Id,
			ShippingFee:              data.ShippingFee,
			TotalProductPrice:        data.TotalProductsPrice,
			MerchantTotalDiscount:    data.MerchantTotalDiscount,
			MarketplaceTotalDiscount: data.MarketPlaceTotalDiscount,
			Transaction:              model.Transaction{},
			Courier:                  model.Courier{},
			User:                     model.User{},
		}

		orderModels = append(orderModels, temp)
	}

	err := o.orderRepository.CreateOrder(user, orderModels, payload.Products, listOfProductVariant)
	if err != nil {
		return nil, err
	}

	resp := dto.StatusMessageResponse{Status: true}
	return &resp, nil
}

func (o orderService) CreateOrderSLP(sign string) (*dto.StatusMessageResponse, error) {
	redisHelperSec := helper.NewRedisClient(7)
	keySec := redisHelperSec.GenerateKey("s", sign, "")
	payloadSec, err := redisHelperSec.GetData(keySec)
	if err != nil {
		return nil, httperror.BadRequestError("Invalid Signature", "INVALID_SIGNATURE")
	}

	defer redisHelperSec.Delete(keySec)

	order := new(dto.CreateOrderRequestSLP)
	json.Unmarshal([]byte(payloadSec), order)

	//verify product
	var listOfProductVariant []uint64
	for _, data := range order.Products {
		listOfProductVariant = append(listOfProductVariant, data.VariantID)
	}

	numOfRecords := o.orderRepository.CountUserCheckoutItem(listOfProductVariant, order.User.Id)
	if numOfRecords != int64(len(listOfProductVariant)) {
		return nil, httperror.BadRequestError("Wrong Product Variant", "WRONG_PRODUCT_VARIANT")
	}

	var orderModels []model.Order
	for i, data := range order.Orders {
		now := time.Now().Format("2006-01-02-150405")

		invoice := "INV/?/TRC-?/MRCH-?/USR-?/?"
		invoice = strings.Replace(invoice, "?", now, 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(order.TransactionId)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(data.MerchantId)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(int(order.User.Id)), 1)
		invoice = strings.Replace(invoice, "?", strconv.Itoa(i), 1)

		temp := model.Order{
			Invoice:                  invoice,
			UserTotalPayment:         data.UserTotalPayment,
			MerchantIncome:           data.MerchantIncome,
			MerchantVoucherId:        data.MerchantVoucherId,
			MarketplaceVoucherId:     data.MarketplaceVoucherId,
			DstAddressId:             data.DstAddressId,
			TransactionId:            order.TransactionId,
			OrderStatusId:            1,
			WithdrawStatus:           false,
			MerchantId:               data.MerchantId,
			CourierId:                data.CourierId,
			UserId:                   order.User.Id,
			ShippingFee:              data.ShippingFee,
			TotalProductPrice:        data.TotalProductsPrice,
			MerchantTotalDiscount:    data.MerchantTotalDiscount,
			MarketplaceTotalDiscount: data.MarketPlaceTotalDiscount,
			Transaction:              model.Transaction{},
			Courier:                  model.Courier{},
			User:                     model.User{},
		}
		orderModels = append(orderModels, temp)
	}

	err = o.orderRepository.CreateOrder(order.User, orderModels, order.Products, listOfProductVariant)
	if err != nil {
		return nil, err
	}

	resp := dto.StatusMessageResponse{Status: true}
	return &resp, nil
}

func (o orderService) GetMerchantOrder(merchantId uint64, keyword string, filters *helper.Filters) ([]model.SellerOrder, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)

	orders, metadata, err := o.orderRepository.GetMerchantOrders(merchantId, keyword, filters)
	if err != nil {
		return nil, nil, err
	}

	return orders, metadata, nil
}

func (o orderService) GetUserOrders(userId uint64, keyword string, filters *helper.Filters) ([]dto.UserOrderRes, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)

	orders, metadata, err := o.orderRepository.GetUserOrders(userId, keyword, filters)
	if err != nil {
		return nil, nil, err
	}

	var dtoRes []dto.UserOrderRes
	for _, order := range orders {
		dtoRes = append(dtoRes, dto.UserOrderRes{
			Id:               order.Id,
			Invoice:          order.Invoice,
			UserTotalPayment: order.UserTotalPayment,
			TransactionId:    order.TransactionId,
			Status:           order.Status,
			Merchant: dto.OrderMerchantRes{
				Id:               order.MerchantId,
				Name:             order.MerchantName,
				Description:      order.MerchantDescription,
				PicUrl:           order.MerchantPicUrl,
				MerchantUsername: order.MerchantUsername,
			},
			CourierName:    order.CourierName,
			OrderDetail:    order.OrderDetail,
			OrderCreatedAt: order.OrderCreatedAt,
		})
	}

	return dtoRes, metadata, nil
}

func (o orderService) GetOrderDetails(orderId, merchantId uint64) (*dto.OrderHighlightResponse, error) {
	order, err := o.orderRepository.GetOrderDetail(orderId, merchantId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidMerchantProduct):
			return nil, httperror.BadRequestError(err.Error(), "INVALID_MERCHANT_ORDER")
		default:
			return nil, err
		}
	}

	return order, nil
}
func (o orderService) GetMerchantOrderSummary(merchantId uint64) ([]model.SellerOrder, error) {
	orders, err := o.orderRepository.GetMerchantOrderSummary(merchantId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func setAtScriptUpdateDelivered(orderId uint64) error {
	command := fmt.Sprintf("at now+10 minutes <<< 'go run scripts/at/update_order_status/main.go %d'", orderId)
	return exec.Command("bash", "-c", command).Run()
}

func (o orderService) UpwardOrderStatus(orderId uint64, user model.User) (*dto.StatusMessageResponse, error) {
	orderModel, err := o.orderRepository.GetCurrentOrder(orderId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	resp := dto.StatusMessageResponse{Status: false}
	if orderModel.Merchant.UserId != user.Id {
		return &resp, httperror.ForbiddenError()
	}
	switch orderModel.OrderStatus.Id {
	case 1:
		orderModel.OrderStatusId = 2
	case 2:
		orderModel.OrderStatusId = 3
	default:
		return &resp, httperror.BadRequestError("Maximum Status Reached", "MAX_STATUS_REACHED")
	}
	err = o.orderRepository.UpdateCurrentOrderId(orderModel)
	if err != nil {
		return nil, err
	}
	resp.Status = true

	if orderModel.OrderStatusId == 3 {
		err := setAtScriptUpdateDelivered(orderId)
		if err != nil {
			return nil, err
		}
	}
	return &resp, nil
}

func (o orderService) CompleteOrder(orderId uint64, user *model.User) error {
	order, err := o.orderRepository.GetCurrentOrder(orderId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return httperror.BadRequestError("no row", "NO_ROWS")
		default:
			return httperror.InternalServerError(err.Error())
		}
	}
	if order.UserId != user.Id {
		return httperror.ForbiddenError()
	}
	if order.OrderStatusId != 4 {
		return httperror.BadRequestError("order status is not valid", "INVALID_ORDER_STATUS")
	}

	order.OrderStatusId = 6
	err = o.orderRepository.UpdateCurrentOrderId(order)
	return err
}
