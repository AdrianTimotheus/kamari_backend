package service

import (
	"errors"
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
	GetTopSellerProduct() ([]dto.Product, error)
	GetProductDetail(productId uint64) (*dto.ProductDetailResponse, error)
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

	cityMapper := make(map[string]bool)
	cities := make([]string, 0)

	categoryMapper := make(map[string]uint64)

	productResponse := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		cityMapper[product.City] = true
		_, ok := categoryMapper[product.Category.Name]
		if !ok {
			categoryMapper[product.Category.Name] = product.Category.Id
		}

		productVar, err := p.productVariantRepository.GetProductVariants(product.Id)
		if err != nil {
			return nil, nil, httperror.InternalServerError(err.Error())
		}

		count, rating, err := p.productRepository.GetProductReviewScoreDecimal(product.Id)
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
			ReviewCount:      count,
			MerchantUsername: product.Seller,
			FavouriteCount:   product.FavouriteCount,
			MerchantId:       product.MerchantId,
			Availability:     product.Availability,
			CategoryId:       product.CategoryId,
			Category:         product.Category,
		}
		productResponse[i] = p
	}

	for c := range cityMapper {
		cities = append(cities, c)
	}

	for cat, id := range categoryMapper {
		metadata.Categories = append(metadata.Categories, helper.NewCategoryMetadata(id, cat))
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
	metadata.Cities = cities

	return productResponse, metadata, nil
}

func (p productService) GetTopSellerProduct() ([]dto.Product, error) {
	productModels, err := p.productRepository.GetTopSellProduct()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	var dtoResponse []dto.Product
	for _, data := range productModels {

		temp := dto.Product{
			ID:           data.Id,
			Title:        data.Name,
			ThumbnailURL: data.ThumbnailUrl,
			MerchantName: data.Merchant.Name,
			ItemSold:     data.UnitSold,
		}

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
		count, productScore, err := p.productRepository.GetProductReviewScore(data.ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		dtoResponse[i].ReviewScore = productScore
		dtoResponse[i].ReviewCount = count
	}
	return dtoResponse, nil
}

func (p productService) GetProductDetail(productId uint64) (*dto.ProductDetailResponse, error) {
	product, err := p.productRepository.GetProductDetails(productId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, err
		}
	}

	count, rating, err := p.productRepository.GetProductReviewScoreDecimal(product.Id)
	if err != nil {
		return nil, err
	}

	var stock int
	for _, variant := range product.ProductVariants {
		stock += variant.Stock
	}

	priceRange := helper.CountPriceRange(product.ProductVariants)
	productResponse := &dto.ProductDetailResponse{
		Id:                      product.Id,
		Title:                   product.Name,
		Description:             product.Description,
		MinimumOrder:            product.MinimumOrder,
		ThumbnailUrl:            product.ThumbnailUrl,
		IsHazardous:             product.IsHazardous,
		Weight:                  product.Weight,
		ViewCount:               product.ViewCount,
		UnitSold:                product.UnitSold,
		FavouriteCount:          product.FavouriteCount,
		Stock:                   stock,
		Availability:            product.Availability,
		SellerCategoryId:        product.SellerCategoryId,
		CategoryId:              product.CategoryId,
		ComplementaryCategoryId: product.ComplementaryCategoryId,
		MerchantID:              product.MerchantID,
		Category:                product.Category,
		Subcategory:             product.Subcategory,
		ProductPhotos:           product.ProductPhotos,
		ProductVariants:         product.ProductVariants,
		ReviewCount:             count,
		ReviewScore:             rating,
		PriceRange:              priceRange,
	}

	return productResponse, nil
}
