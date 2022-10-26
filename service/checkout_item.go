package service

import (
	"errors"
	"fmt"
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
	merchantRepository         repository.MerchantRepository
	checkoutRepository         repository.CheckoutItemRepository
	productRepository          repository.ProductRepository
	productPromoItemRepository repository.ProductPromoItemRepository
}

type CISConfig struct {
	CheckoutRepository         repository.CheckoutItemRepository
	ProductRepository          repository.ProductRepository
	ProductPromoItemRepository repository.ProductPromoItemRepository
	MerchantRepository         repository.MerchantRepository
}

func NewCheckoutItemService(c CISConfig) CheckoutItemService {
	return &checkoutItemService{
		checkoutRepository:         c.CheckoutRepository,
		productRepository:          c.ProductRepository,
		productPromoItemRepository: c.ProductPromoItemRepository,
		merchantRepository:         c.MerchantRepository,
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
			merchantModel, err := c.merchantRepository.GetMerchantDetailById(temp.Merchant.Id)
			if err != nil {
				switch {
				case errors.Is(err, gorm.ErrRecordNotFound):
				default:
					return nil, httperror.InternalServerError(err.Error())
				}
			}
			temp.Merchant.CityId = merchantModel.UserAddress.CityId
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
			RealPrice:    data.ProductVariant.Price,
			Price:        data.ProductVariant.Price,
			TotalPrice:   data.ProductVariant.Price * uint64(data.Quantity),
			ThumbnailURL: data.ProductVariant.PicUrl,
			MainVariant: dto.MainVariant{
				Name: data.ProductVariant.MainVariant,
				Item: data.ProductVariant.MainVariantItem,
			},
			SubVariant: dto.SubVariant{
				Name: data.ProductVariant.SubVariant,
				Item: data.ProductVariant.SubVariantItem,
			},
			Quantity:    data.Quantity,
			ItemWeight:  productData[data.ProductVariant.ProductID].Weight,
			TotalWeight: uint64(data.Quantity) * productData[data.ProductVariant.ProductID].Weight,
		}
		dtoResponse[sliceIndex].Products = append(dtoResponse[sliceIndex].Products, tempProduct)
	}
	for i, merchant := range dtoResponse {
		var listOfProductId []uint64
		for j, product := range merchant.Products {
			listOfProductId = append(listOfProductId, product.ProductId)
			promoItemModel, err := c.productPromoItemRepository.GetProductPromoItemByProductId(product.ProductId)
			if err != nil {
				return nil, err
			}
			if promoItemModel != nil {
				if promoItemModel.MaximumItem >= uint64(product.Quantity) {
					switch promoItemModel.PromoTypeId {
					case 1:
						if dtoResponse[i].Products[j].Price = product.Price - promoItemModel.DiscountAmount; dtoResponse[i].Products[j].Price < 0 {
							dtoResponse[i].Products[j].Price = 0
							dtoResponse[i].Products[j].TotalPrice = 0
						} else {
							fmt.Println(dtoResponse[i].Products[j].Price)
							fmt.Println(dtoResponse[i].Products[j].TotalPrice)
							dtoResponse[i].Products[j].TotalPrice = uint64(product.Quantity) * dtoResponse[i].Products[j].Price
							fmt.Println(dtoResponse[i].Products[j].TotalPrice)
						}
					case 2:
						if promoItemModel.DiscountAmount == 100 {
							dtoResponse[i].Products[j].Price = 0
							dtoResponse[i].Products[j].TotalPrice = 0
						} else {
							dtoResponse[i].Products[j].Price = product.Price - ((product.Price * promoItemModel.DiscountAmount) / 100)
							dtoResponse[i].Products[j].TotalPrice = dtoResponse[i].Products[j].Price * uint64(product.Quantity)
						}
					}
				}
			}
		}
		courierModels, err := c.checkoutRepository.GetSupportedCourierForMerchant(merchant.Merchant.Id, listOfProductId)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):

			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		for _, courier := range courierModels {
			temp := dto.SupportedCourier{
				Id:   courier.Id,
				Name: courier.Name,
			}
			dtoResponse[i].Couriers = append(dtoResponse[i].Couriers, temp)
		}
	}
	return dtoResponse, nil
}
