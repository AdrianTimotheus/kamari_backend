package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/service"
)

func TestNewMerchantService(t *testing.T) {
	mockRepo := &mocks.MerchantRepository{}
	userService := service.NewMerchantService(service.MSConfig{
		MerchantRepository: mockRepo,
	})

	assert.NotNil(t, userService)
}

func TestMerchantService_GetMerchantDetail(t *testing.T) {
	t.Run("should thrown error if merchant not found", func(t *testing.T) {
		mockRepo := &mocks.MerchantRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository: mockRepo,
		})
		mockRepo.On("GetMerchantDetail", "anonymous").Return(model.Merchant{}, gorm.ErrRecordNotFound)

		_, err := m.GetMerchantDetail("anonymous")

		assert.Error(t, err)
		assert.ErrorAs(t, err, &gorm.ErrRecordNotFound)
	})

	t.Run("should thrown error if merchant not found", func(t *testing.T) {
		var er httperror.AppError
		mockRepo := &mocks.MerchantRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository: mockRepo,
		})
		mockRepo.On("GetMerchantDetail", "anonymous").Return(model.Merchant{}, httperror.InternalServerError("error"))

		_, err := m.GetMerchantDetail("anonymous")

		assert.Error(t, err)
		assert.ErrorAs(t, err, &er)
	})

	t.Run("should thrown error if failed to fetch merchant location (city)", func(t *testing.T) {
		var er httperror.AppError
		mockRepo := &mocks.MerchantRepository{}
		cRepo := &mocks.CityRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository: mockRepo,
			CityRepository:     cRepo,
		})
		ud := model.UserAddress{CityId: 2}
		mockRepo.On("GetMerchantDetail", "anonymous").Return(model.Merchant{UserAddress: ud}, nil)
		cRepo.On("GetCity", ud.CityId).Return(model.City{}, httperror.InternalServerError("error"))

		merchant, err := m.GetMerchantDetail("anonymous")

		assert.Error(t, err)
		assert.ErrorAs(t, err, &er)
		assert.Nil(t, merchant)
	})

	t.Run("should thrown error if failed to fetch merchant's product", func(t *testing.T) {
		var er httperror.AppError
		mockRepo := &mocks.MerchantRepository{}
		cRepo := &mocks.CityRepository{}
		pRepo := &mocks.ProductRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository: mockRepo,
			CityRepository:     cRepo,
			ProductRepository:  pRepo,
		})
		ud := model.Merchant{
			Id: 1,
			UserAddress: model.UserAddress{
				CityId: 2,
			},
		}
		mockRepo.On("GetMerchantDetail", "anonymous").Return(ud, nil)
		cRepo.On("GetCity", ud.UserAddress.CityId).Return(model.City{}, nil)
		pRepo.On("GetMerchantProduct", ud.Id).Return(nil, httperror.InternalServerError("error"))

		merchant, err := m.GetMerchantDetail("anonymous")

		assert.Error(t, err)
		assert.ErrorAs(t, err, &er)
		assert.Nil(t, merchant)
	})

	t.Run("should thrown error if failed to get merchant's follower count", func(t *testing.T) {
		var er httperror.AppError
		mockRepo := &mocks.MerchantRepository{}
		cRepo := &mocks.CityRepository{}
		pRepo := &mocks.ProductRepository{}
		fRepo := &mocks.MerchantFollowerRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository:         mockRepo,
			CityRepository:             cRepo,
			ProductRepository:          pRepo,
			MerchantFollowerRepository: fRepo,
		})
		ud := model.Merchant{
			Id: 1,
			UserAddress: model.UserAddress{
				CityId: 2,
			},
		}
		mockRepo.On("GetMerchantDetail", "anonymous").Return(ud, nil)
		cRepo.On("GetCity", ud.UserAddress.CityId).Return(model.City{}, nil)
		pRepo.On("GetMerchantProduct", ud.Id).Return(nil, nil)
		fRepo.On("GetMerchantFollowerCount", ud.Id).Return(uint64(0), httperror.InternalServerError("error"))

		merchant, err := m.GetMerchantDetail("anonymous")

		assert.Error(t, err)
		assert.ErrorAs(t, err, &er)
		assert.Nil(t, merchant)
	})

	t.Run("should thrown error if failed to get merchant's follower count", func(t *testing.T) {
		mockRepo := &mocks.MerchantRepository{}
		cRepo := &mocks.CityRepository{}
		pRepo := &mocks.ProductRepository{}
		fRepo := &mocks.MerchantFollowerRepository{}
		m := service.NewMerchantService(service.MSConfig{
			MerchantRepository:         mockRepo,
			CityRepository:             cRepo,
			ProductRepository:          pRepo,
			MerchantFollowerRepository: fRepo,
		})
		ud := model.Merchant{
			Id: 1,
			UserAddress: model.UserAddress{
				CityId: 2,
			},
		}
		products := []model.Product{{Name: "something"}}
		mockRepo.On("GetMerchantDetail", "anonymous").Return(ud, nil)
		cRepo.On("GetCity", ud.UserAddress.CityId).Return(model.City{}, nil)
		pRepo.On("GetMerchantProduct", ud.Id).Return(products, nil)
		fRepo.On("GetMerchantFollowerCount", ud.Id).Return(uint64(12), nil)

		merchant, err := m.GetMerchantDetail("anonymous")

		assert.Nil(t, err)
		assert.NotNil(t, merchant)
	})
}
