package service

import (
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/repository"
	"time"
)

type MarketplaceService interface {
	GetPromoBanner() ([]dto.GetPromoBanner, error)
}

type marketplaceService struct {
	marketplaceRepository repository.MarketplaceRepository
}

type MpSConfig struct {
	MarketplaceRepository repository.MarketplaceRepository
}

const (
	YYYYMMDD = "2006-01-02"
)

func NewMarketplaceService(config MpSConfig) MarketplaceService {
	return &marketplaceService{
		marketplaceRepository: config.MarketplaceRepository,
	}
}

func (m marketplaceService) GetPromoBanner() ([]dto.GetPromoBanner, error) {
	dateNow := time.Now()
	dateString := dateNow.Format(YYYYMMDD)
	results, err := m.marketplaceRepository.GetPromoBanner(dateString)
	if err != nil {
		return nil, httperror.InternalServerError(err.Error())
	}

	var dtoResponse []dto.GetPromoBanner
	for _, data := range results {
		temp := dto.GetPromoBanner{
			ID:     data.ID,
			PicUrl: data.PicUrl,
		}
		dtoResponse = append(dtoResponse, temp)
	}
	return dtoResponse, nil
}
