package service_test

import (
	"github.com/stretchr/testify/assert"
	"kamari.com/httperror"
	"kamari.com/mocks"
	"kamari.com/model"
	"kamari.com/service"
	"testing"
)

func TestNewProvinceService(t *testing.T) {
	mockRepo := &mocks.ProvinceRepository{}
	provinceService := service.NewProvinceService(service.PSConfig{
		ProvinceRepository: mockRepo,
	})

	assert.NotNil(t, provinceService)
}

func TestProvinceService_GetAllProvinces(t *testing.T) {
	t.Run("Should return all provinces by default", func(t *testing.T) {
		mockRepo := &mocks.ProvinceRepository{}
		s := service.NewProvinceService(service.PSConfig{
			ProvinceRepository: mockRepo,
		})

		data := []model.Province{
			model.Province{
				Id:   1,
				Name: "DKI Jakarta",
			},
			model.Province{
				Id:   2,
				Name: "Banten",
			},
		}

		mockRepo.On("GetAllProvinces").Return(data, nil)

		provinceData, err := s.GetAllProvinces()

		assert.Nil(t, err)
		assert.Equal(t, data, provinceData)
	})

	t.Run("Should return error if something went wrong", func(t *testing.T) {
		mockRepo := &mocks.ProvinceRepository{}
		s := service.NewProvinceService(service.PSConfig{
			ProvinceRepository: mockRepo,
		})

		mockRepo.On("GetAllProvinces").Return(nil, httperror.InternalServerError("Internal server error"))

		provinceData, err := s.GetAllProvinces()

		assert.Nil(t, provinceData)
		assert.Equal(t, err, httperror.InternalServerError("Internal server error"))
	})
}
