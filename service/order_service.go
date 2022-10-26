package service

import (
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"strconv"
	"strings"
	"time"
)

type OrderService interface {
	CreateOrder(user model.User, payload dto.CreateOrderRequest) (*dto.StatusMessageResponse, error)
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

	err := o.orderRepository.CreateOrder(user, orderModels, payload.Products)
	if err != nil {
		return nil, err
	}

	resp := dto.StatusMessageResponse{Status: true}
	return &resp, nil
}
