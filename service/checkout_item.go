package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type CheckoutItemService interface {
	CreateCheckoutItem(model.User, []dto.PostCheckoutRequest) (dto.StatusMessageResponse, error)
	GetCheckoutItemByUser(model.User) ([]dto.CheckoutItemResponse, error)
}

type checkoutItemService struct {
	checkoutRepository repository.CheckoutItemRepository
	productRepository  repository.ProductRepository
}

type CISConfig struct {
	CheckoutRepository repository.CheckoutItemRepository
	ProductRepository  repository.ProductRepository
}

func NewCheckoutItemService(c CISConfig) CheckoutItemService {
	return &checkoutItemService{
		checkoutRepository: c.CheckoutRepository,
		productRepository:  c.ProductRepository,
	}
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

func (c checkoutItemService) GetCheckoutItemByUser(user model.User) ([]dto.CheckoutItemResponse, error) {
	models, err := c.checkoutRepository.GetCheckoutItemByUser(user.Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_DATA")
		default:
			return nil, err
		}
	}
	//index for dto based on merchant
	dtoIndex := make(map[string]int)
	productData := make(map[uint64]model.Product)
	var dtoResponse []dto.CheckoutItemResponse
	initIndex := 0
	for _, data := range models {
		_, ok := dtoIndex[data.Merchant.Name]
		//merchantData
		if !ok {
			dtoIndex[data.Merchant.Name] = initIndex
			initIndex++
			temp := dto.CheckoutItemResponse{
				Merchant: dto.MerchantDetailResponse{
					Id:             data.Merchant.Id,
					Name:           data.Merchant.Name,
					ProfilePicture: data.Merchant.PicUrl,
				},
			}
			dtoResponse = append(dtoResponse, temp)
		}
		//productData
		_, ok = productData[data.ProductVariant.ProductID]
		if !ok {
			product, err := c.productRepository.GetProductPlainDetail(data.ProductVariant.ProductID)
			if err != nil {
				return nil, err
			}
			productData[data.ProductVariant.ProductID] = *product
		}
		//assign product data
		sliceIndex, _ := dtoIndex[data.Merchant.Name]
		tempProduct := dto.CheckoutProduct{
			ProductId:    data.ProductVariant.ProductID,
			VariantID:    data.ProductVariantId,
			ProductName:  productData[data.ProductVariant.ProductID].Name,
			Price:        data.ProductVariant.Price,
			ThumbnailURL: data.ProductVariant.PicUrl,
			MainVariant: dto.MainVariant{
				Name: data.ProductVariant.MainVariant,
				Item: data.ProductVariant.MainVariantItem,
			},
			SubVariant: dto.SubVariant{
				Name: data.ProductVariant.SubVariant,
				Item: data.ProductVariant.SubVariantItem,
			},
			Quantity: data.Quantity,
		}
		dtoResponse[sliceIndex].Products = append(dtoResponse[sliceIndex].Products, tempProduct)
	}
	return dtoResponse, nil
}
