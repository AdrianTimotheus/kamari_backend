package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type ProductService interface {
	GetProducts(keyword string, filters *helper.Filters) ([]dto.ProductResponse, *helper.Metadata, error)
	GetTopSellerProduct() ([]dto.DtoProducts, error)
}

type productService struct {
	productRepository        repository.ProductRepository
	productVariantRepository repository.ProductVariantRepository
	orderDetailsRepository   repository.OrderDetailsRepository
	cityRepository           repository.CityRepository
	merchantRepository       repository.MerchantRepository
}

type PConfig struct {
	ProductRepository        repository.ProductRepository
	ProductVariantRepository repository.ProductVariantRepository
	OrderDetailsRepository   repository.OrderDetailsRepository
	CityRepository           repository.CityRepository
	MerchantRepository       repository.MerchantRepository
}

func NewProductService(c PConfig) ProductService {
	return productService{
		productRepository:        c.ProductRepository,
		productVariantRepository: c.ProductVariantRepository,
		orderDetailsRepository:   c.OrderDetailsRepository,
		cityRepository:           c.CityRepository,
		merchantRepository:       c.MerchantRepository,
	}
}

func (p productService) GetProducts(keyword string, filters *helper.Filters) ([]dto.ProductResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)
	filters.SortDirection = strings.ToLower(helper.ReadString(filters.SortDirection, "ASC"))
	filters.Sort = helper.ReadString(filters.Sort, "id")
	filters.SortBy = strings.ToLower(filters.SortBy)

	products, metadata, err := p.productRepository.GetAllProducts(keyword, filters)
	if err != nil {
		return nil, nil, httperror.InternalServerError(err.Error())
	}

	productResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productVar, err := p.productVariantRepository.GetProductVariants(product.Id)
		if err != nil {
			return nil, nil, httperror.InternalServerError(err.Error())
		}

		rating, err := p.productRepository.GetProductReviewScore(product.Id)
		if err != nil {
			return nil, nil, httperror.InternalServerError(err.Error())
		}

		lowestProdPrice := helper.GetMinProductVarPrice(productVar)
		p := dto.ProductResponse{
			Id:               product.Id,
			Title:            product.Name,
			Price:            lowestProdPrice,
			ItemSold:         product.UnitSold,
			City:             product.City,
			ThumbnailUrl:     product.ThumbnailUrl,
			ReviewScore:      rating,
			MerchantUsername: product.Seller,
			FavouriteCount:   product.FavouriteCount,
			MerchantId:       product.MerchantId,
			Availability:     product.Availability,
			CategoryId:       product.CategoryId,
			Category:         product.Category,
		}
		productResponse[i] = p
	}

	if filters.SortBy == "price" {
		switch filters.SortDirection {
		case "asc":
			sort.Slice(productResponse, func(i, j int) bool {
				return productResponse[i].Price < productResponse[j].Price
			})
		case "desc":
			sort.Slice(productResponse, func(i, j int) bool {
				return productResponse[i].Price > productResponse[j].Price
			})
		}
	}

	if filters.MinPrice != nil {
		productResponse = helper.FilterProductPrice(productResponse, "minPrice", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		productResponse = helper.FilterProductPrice(productResponse, "maxPrice", *filters.MaxPrice)
	}
	if filters.Rating != nil {
		productResponse = helper.FilterProductByRating(productResponse, *filters.Rating)
	}
	if filters.Location != "" {
		productResponse = helper.FilterBySellerLocation(productResponse, filters.Location)
	}

	metadata.TotalRecords = len(productResponse)

	return productResponse, metadata, nil
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
