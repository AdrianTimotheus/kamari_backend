package service

import (
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type CheckoutItemService interface {
	CreateCheckoutItem(model.User, []dto.PostCheckoutRequest) (dto.StatusMessageResponse, error)
}

type checkoutItemService struct {
	checkoutRepository repository.CheckoutItemRepository
}

type CISConfig struct {
	CheckoutRepository repository.CheckoutItemRepository
}

func NewCheckoutItemService(c CISConfig) CheckoutItemService {
	return &checkoutItemService{checkoutRepository: c.CheckoutRepository}
}

func (c checkoutItemService) CreateCheckoutItem(userModel model.User, requests []dto.PostCheckoutRequest) (dto.StatusMessageResponse, error) {
	var models []model.CheckoutItem
	for _, data := range requests {
		tempModel := model.CheckoutItem{
			ProductVariantId: data.ProductVariantId,
			MerchantId:       data.MerchantId,
			UserId:           userModel.Id,
			Quantity:         data.Quantity,
		}
		models = append(models, tempModel)
	}
	err := c.checkoutRepository.CreateCheckoutRecords(models)
	if err != nil {
		return dto.StatusMessageResponse{Status: false}, httperror.InternalServerError("Internal Server Error")
	}
	return dto.StatusMessageResponse{Status: true}, nil
}
