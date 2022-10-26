package service

import (
	"errors"
	"kamari.com/model"
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
	GetSimilarProduct(payload dto.GetSimilarProductRequest) ([]dto.Product, error)
	AddNewProduct(payload dto.AddProductReq, user model.User) error
	EditProduct(payload dto.EditProductReq, user model.User, id uint64) error
	DeleteProduct(id uint64, user model.User) error
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
		count, productScore, err := p.productRepository.GetProductReviewScoreDecimal(data.ID)
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

	//for i := range product.ProductVariants {
	//	if product.ProductVariants[i].ProductPromoID != 0 {
	//		switch product.ProductVariants[i].ProductPromo.PromoTypeId {
	//		case 1:
	//			product.ProductVariants[i].PromoPrice = product.ProductVariants[i].Price - product.ProductVariants[i].ProductPromo.DiscountAmount
	//		case 2:
	//			p := product.ProductVariants[i].Price
	//			discountAmount := p * product.ProductVariants[i].ProductPromo.DiscountAmount / 100
	//			product.ProductVariants[i].PromoPrice = p - discountAmount
	//		}
	//	}
	//}

	// product promo
	var promo *dto.PromoItems
	promotion, err := p.productRepository.GetProductPromotion(productId)
	if err == nil {
		promo = &dto.PromoItems{
			Id:             promotion.Id,
			ProductId:      promotion.ProductId,
			PromoTypeId:    promotion.PromoTypeId,
			MaximumItem:    promotion.MaximumItem,
			DiscountAmount: promotion.DiscountAmount,
		}
	}

	priceRange := helper.CountPriceRange(product.ProductVariants)
	productResponse := &dto.ProductDetailResponse{
		Id:                      product.Id,
		Title:                   product.Name,
		Description:             product.Description,
		MinimumOrder:            product.MinimumOrder,
		MaximumOrder:            product.MaximumOrder,
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
		MerchantUsername:        product.Merchant.User.Username,
		Category:                product.Category,
		Subcategory:             product.Subcategory,
		ProductPhotos:           product.ProductPhotos,
		ProductVariants:         product.ProductVariants,
		ReviewCount:             count,
		ReviewScore:             rating,
		PriceRange:              priceRange,
		Promo:                   promo,
	}

	return productResponse, nil
}

func (p productService) GetSimilarProduct(payload dto.GetSimilarProductRequest) ([]dto.Product, error) {
	currentProductModel, err := p.productRepository.GetProductDetails(payload.CurrentProductId)
	if err != nil {
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
	}

	models, err := p.productRepository.GetSimilarProduct(payload.Keyword, currentProductModel.CategoryId, payload.CurrentProductId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	var dtoResponse []dto.Product
	for _, data := range models {
		temp := dto.Product{
			ID:           data.Id,
			Title:        data.Name,
			ThumbnailURL: data.ThumbnailUrl,
			MerchantName: data.Merchant.Name,
			ItemSold:     data.UnitSold,
			//ReviewScore:  0,
			//ReviewCount:  0,
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
		pvModel, err := p.productRepository.GetProductFirstVariant(temp.ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}

		//product score
		count, productScore, err := p.productRepository.GetProductReviewScoreDecimal(temp.ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		temp.Price = pvModel.Price
		temp.ProvinceName = cityModel.Province.Name
		temp.CityName = cityModel.Name
		temp.ReviewScore = productScore
		temp.ReviewCount = count
		dtoResponse = append(dtoResponse, temp)
	}

	return dtoResponse, nil
}

func (p productService) AddNewProduct(payload dto.AddProductReq, user model.User) error {
	newProduct := model.Product{
		Name:                    payload.Name,
		Description:             payload.Description,
		MinimumOrder:            payload.MinimumOrder,
		MaximumOrder:            payload.MaximumOrder,
		ThumbnailUrl:            payload.ThumbnailUrl,
		IsHazardous:             *payload.IsHazardous,
		Weight:                  payload.Weight,
		CategoryId:              payload.CategoryId,
		SubcategoryId:           payload.SubcategoryId,
		ComplementaryCategoryId: payload.ComplementaryCategoryId,
		ProductVariants:         payload.ProductVariants,
	}

	err := p.productRepository.AddNewProduct(newProduct, user)
	if err != nil {
		return err
	}
	return nil
}

func (p productService) EditProduct(payload dto.EditProductReq, user model.User, id uint64) error {
	newProduct := model.Product{
		Name:                    payload.Name,
		Description:             payload.Description,
		MinimumOrder:            payload.MinimumOrder,
		MaximumOrder:            payload.MaximumOrder,
		ThumbnailUrl:            payload.ThumbnailUrl,
		IsHazardous:             payload.IsHazardous,
		Weight:                  payload.Weight,
		CategoryId:              payload.CategoryId,
		SellerCategoryId:        payload.SellerCategoryId,
		SubcategoryId:           payload.SubcategoryId,
		ComplementaryCategoryId: payload.ComplementaryCategoryId,
		ProductVariants:         payload.ProductVariants,
	}
	err := p.productRepository.DeleteProduct(id, user)
	if err != nil {
		return err
	}
	err = p.productRepository.AddNewProduct(newProduct, user)
	if err != nil {
		return err
	}
	return nil
}

func (p productService) DeleteProduct(id uint64, user model.User) error {
	err := p.productRepository.DeleteProduct(id, user)
	if err != nil {
		return err
	}
	return nil
}
