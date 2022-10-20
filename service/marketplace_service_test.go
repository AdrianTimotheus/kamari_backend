package service_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/service"
	"testing"
	"time"
)

const (
	YYYYMMDD = "2006-01-02"
)

func TestNewMarketplaceService(t *testing.T) {
	t.Run("Should return object on success", func(t *testing.T) {
		mockRepo := &mocks.MarketplaceRepository{}
		marketplaceService := service.NewMarketplaceService(service.MpSConfig{
			MarketplaceRepository: mockRepo,
		})
		assert.NotNil(t, marketplaceService)
	})
}

func TestMarketplaceService_GetPromoBanner(t *testing.T) {
	t.Run("should return DTO response on success", func(t *testing.T) {
		dateNow := time.Now()
		dateString := dateNow.Format(YYYYMMDD)
		mockRepo := &mocks.MarketplaceRepository{}
		returnModel := []model.MarketplaceVoucher{
			model.MarketplaceVoucher{
				ID:             1,
				Title:          "test",
				Code:           "test",
				PromoTypeID:    1,
				DiscountAmount: 100,
				StartDate:      time.Now(),
				EndDate:        time.Now(),
				StockLeft:      10,
				PicUrl:         "test",
			},
			model.MarketplaceVoucher{
				ID:             2,
				Title:          "test",
				Code:           "test",
				PromoTypeID:    1,
				DiscountAmount: 100,
				StartDate:      time.Now(),
				EndDate:        time.Now(),
				StockLeft:      10,
				PicUrl:         "test",
			},
		}

		mockRepo.On("GetPromoBanner", dateString).Return(returnModel, nil)

		services := service.NewMarketplaceService(service.MpSConfig{
			MarketplaceRepository: mockRepo,
		})
		dto, err := services.GetPromoBanner()

		assert.NotNil(t, dto)
		assert.Equal(t, 2, len(dto))
		assert.Nil(t, err)
	})

	t.Run("should return error on failed fetch", func(t *testing.T) {
		dateNow := time.Now()
		dateString := dateNow.Format(YYYYMMDD)
		mockRepo := &mocks.MarketplaceRepository{}
		mockRepo.On("GetPromoBanner", dateString).Return(nil, errors.New("test"))

		services := service.NewMarketplaceService(service.MpSConfig{
			MarketplaceRepository: mockRepo,
		})
		dto, err := services.GetPromoBanner()

		assert.Nil(t, dto)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "test")
	})
}
