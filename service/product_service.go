package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type ProductService interface {
	GetTopSellerProduct() ([]dto.DtoProducts, error)
}

type productService struct {
	productRepository  repository.ProductRepository
	merchantRepository repository.MerchantRepository
	cityRepository     repository.CityRepository
}

type PConfig struct {
	ProductRepository  repository.ProductRepository
	MerchantRepository repository.MerchantRepository
	CityRepository     repository.CityRepository
}

func NewProductService(c PConfig) ProductService {
	return productService{
		productRepository:  c.ProductRepository,
		merchantRepository: c.MerchantRepository,
		cityRepository:     c.CityRepository,
	}
}

func (p productService) GetTopSellerProduct() ([]dto.DtoProducts, error) {
	productModels, err := p.productRepository.GetTopSellProduct()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	var dtoResponse []dto.DtoProducts
	for _, data := range productModels {

		temp := dto.DtoProducts{
			ID:           data.Id,
			Title:        data.Name,
			ThumbnailURL: data.ThumbnailUrl,
			//ProvinceName: "",
			//CityName:     "",
			MerchantName: data.Merchant.Name,
			ItemSold:     data.UnitSold,
			//ReviewScore:  0,
		}
		fmt.Println("title", data.Name)
		fmt.Println("Merchant ID", data.MerchantID)

		merchantModel, err := p.merchantRepository.GetMerchantDetailById(data.MerchantID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}

		cityModel, err := p.cityRepository.GetCity(merchantModel.User.CityId)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		temp.CityName = cityModel.Name
		temp.ProvinceName = cityModel.Province.Name
		dtoResponse = append(dtoResponse, temp)
	}
	for i, data := range dtoResponse {
		//first variant product for price
		pvModel, err := p.productRepository.GetProductFirstVariant(data.ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		dtoResponse[i].Price = pvModel.Price

		//product score
		productScore, err := p.productRepository.GetProductReviewScore(data.ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		dtoResponse[i].ReviewScore = productScore
	}
	return dtoResponse, nil
}
