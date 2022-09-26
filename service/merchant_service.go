package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type MerchantService interface {
	GetMerchantDetail(username string) (*dto.MerchantDetailResponse, error)
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
		Name:              md.Name,
		Description:       md.Description,
		ProfilePicture:    md.User.PicUrl,
		ShopLocation:      ml.Name,
		FollowerCount:     mf,
		ProductCount:      uint64(len(mp)),
		MerchantCreatedAt: md.CreatedAt,
	}

	return merchant, nil
}
