package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"sort"
)

type MerchantService interface {
	GetMerchantDetail(username string) (*dto.MerchantDetailResponse, error)
	GetMerchantProducts(username string, payload dto.GetMerchantProductRequest) (*dto.GetMerchantProductsResponse, error)
	GetSellerAdminProducts(username string, payload dto.GetMerchantProductRequest) (*dto.GetSellerAdminProductsResponse, error)
	RegisterMerchant(user model.User, req *dto.MerchantRegisterRequest) error
}

type merchantService struct {
	merchantRepository         repository.MerchantRepository
	cityRepository             repository.CityRepository
	productRepository          repository.ProductRepository
	merchantFollowerRepository repository.MerchantFollowerRepository
}

type MSConfig struct {
	MerchantRepository         repository.MerchantRepository
	CityRepository             repository.CityRepository
	ProductRepository          repository.ProductRepository
	MerchantFollowerRepository repository.MerchantFollowerRepository
}

func NewMerchantService(c MSConfig) MerchantService {
	return merchantService{
		merchantRepository:         c.MerchantRepository,
		cityRepository:             c.CityRepository,
		productRepository:          c.ProductRepository,
		merchantFollowerRepository: c.MerchantFollowerRepository,
	}
}

func (m merchantService) GetMerchantDetail(username string) (*dto.MerchantDetailResponse, error) {
	md, err := m.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	ml, err := m.cityRepository.GetCity(md.UserAddress.CityId)
	if err != nil {
		return nil, httperror.InternalServerError(err.Error())
	}

	mp, err := m.productRepository.GetMerchantProduct(md.Id)
	if err != nil {
		return nil, httperror.InternalServerError(err.Error())
	}

	mf, err := m.merchantFollowerRepository.GetMerchantFollowerCount(md.Id)
	if err != nil {
		return nil, httperror.InternalServerError(err.Error())
	}

	merchant := &dto.MerchantDetailResponse{
		Id:                md.Id,
		Name:              md.Name,
		Username:          md.User.Username,
		Description:       md.Description,
		ProfilePicture:    md.User.PicUrl,
		ShopLocation:      ml.Name,
		FollowerCount:     mf,
		ProductCount:      uint64(len(mp)),
		MerchantCreatedAt: md.CreatedAt,
	}

	return merchant, nil
}

func (m merchantService) GetMerchantProducts(username string, payload dto.GetMerchantProductRequest) (*dto.GetMerchantProductsResponse, error) {
	var dtoResp dto.GetMerchantProductsResponse

	merchantModel, err := m.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	cityModel, err := m.cityRepository.GetCity(merchantModel.UserAddress.CityId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	if payload.Page <= 0 {
		payload.Page = 1
	}
	pageOffset := (payload.Page * payload.Limit) - payload.Limit
	var orderBy string
	switch payload.SortBy {
	case "latest":
		orderBy = "created_at DESC"
	case "oldest":
		orderBy = "created_at ASC"
	default:
		orderBy = "unit_sold DESC"
	}
	var models []model.Product
	if payload.Search != "" {
		models, err = m.productRepository.GetSearchProducts(merchantModel.Id, pageOffset, orderBy, payload.Search, payload.Limit, payload.SellerCategory, payload.MarketplaceCategory)
	} else {
		models, err = m.productRepository.GetMerchantProducts(merchantModel.Id, pageOffset, orderBy, payload.Limit, payload.SellerCategory, payload.MarketplaceCategory)
	}

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	var dtoItem []dto.Product
	for _, modelData := range models {
		temp := dto.Product{
			ID:           modelData.Id,
			Title:        modelData.Name,
			ThumbnailURL: modelData.ThumbnailUrl,
			ItemSold:     modelData.UnitSold,
			ProvinceName: cityModel.Province.Name,
			CityName:     cityModel.Name,
			MerchantName: merchantModel.Name,
		}
		dtoItem = append(dtoItem, temp)
	}

	for i, _ := range dtoItem {
		count, score, err := m.productRepository.GetProductReviewScoreDecimal(dtoItem[i].ID)
		if err != nil {
			return nil, httperror.InternalServerError(err.Error())
		}
		dtoItem[i].ReviewScore = score
		dtoItem[i].ReviewCount = count
	}

	for i, _ := range dtoItem {
		model, err := m.productRepository.GetProductFirstVariant(dtoItem[i].ID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
			default:
				return nil, httperror.InternalServerError(err.Error())
			}
		}
		dtoItem[i].Price = model.Price
	}

	switch payload.SortBy {
	case "lowestPrice":
		sort.Slice(dtoItem, func(i, j int) bool {
			return dtoItem[i].Price < dtoItem[j].Price
		})
		for _, v := range dtoItem {
			fmt.Println(v.Price)
		}
	case "highestPrice":
		sort.Slice(dtoItem, func(i, j int) bool {
			return dtoItem[i].Price > dtoItem[j].Price
		})
		for _, v := range dtoItem {
			fmt.Println(v.Price)
		}
	}

	totalRecord, _ := m.productRepository.GetCountMerchantProduct(username)

	if int(totalRecord)%payload.Limit == 0 {
		dtoResp.TotalPage = int(totalRecord) / payload.Limit
	} else {
		dtoResp.TotalPage = (int(totalRecord) / payload.Limit) + 1
	}

	dtoResp.Page = payload.Page
	dtoResp.Products = dtoItem
	return &dtoResp, nil
}

func (m merchantService) GetSellerAdminProducts(username string, payload dto.GetMerchantProductRequest) (*dto.GetSellerAdminProductsResponse, error) {
	var dtoResp dto.GetSellerAdminProductsResponse

	merchantModel, err := m.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	if payload.Page <= 0 {
		payload.Page = 1
	}
	pageOffset := (payload.Page * payload.Limit) - payload.Limit
	var orderBy string
	switch payload.SortBy {
	case "latest":
		orderBy = "created_at DESC"
	case "oldest":
		orderBy = "created_at ASC"
	default:
		orderBy = "unit_sold DESC"
	}
	var models []model.Product
	models, err = m.productRepository.GetSellerAdminProducts(merchantModel.Id, pageOffset, orderBy, payload.Search, payload.Limit, payload.SellerCategory, payload.MarketplaceCategory)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	var dtoItem []dto.SellerAdminProduct
	for _, modelData := range models {
		temp := dto.SellerAdminProduct{
			ID:              modelData.Id,
			Title:           modelData.Name,
			ThumbnailURL:    modelData.ThumbnailUrl,
			ItemSold:        modelData.UnitSold,
			ProductVariants: modelData.ProductVariants,
		}
		dtoItem = append(dtoItem, temp)
	}

	for i, _ := range dtoItem {
		count, score, err := m.productRepository.GetProductReviewScoreDecimal(dtoItem[i].ID)
		if err != nil {
			return nil, httperror.InternalServerError(err.Error())
		}
		dtoItem[i].ReviewScore = score
		dtoItem[i].ReviewCount = count
	}

	totalRecord, _ := m.productRepository.GetCountMerchantProduct(username)

	if int(totalRecord)%payload.Limit == 0 {
		dtoResp.TotalPage = int(totalRecord) / payload.Limit
	} else {
		dtoResp.TotalPage = (int(totalRecord) / payload.Limit) + 1
	}

	dtoResp.Page = payload.Page
	dtoResp.Products = dtoItem
	return &dtoResp, nil
}

func (m merchantService) RegisterMerchant(user model.User, req *dto.MerchantRegisterRequest) error {
	err := m.merchantRepository.RegisterMerchant(user.Id, req)
	if err != nil {
		return err
	}
	return nil
}
