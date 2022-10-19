package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type CourierService interface {
	GetAllCouriers() (*dto.AllCouriersResponse, error)
}

type courierService struct {
	courierRepository repository.CourierRepository
}

type CourierSConfig struct {
	CourierRepository repository.CourierRepository
}

func NewCourierService(c CourierSConfig) CourierService {
	return &courierService{courierRepository: c.CourierRepository}
}

func (c courierService) GetAllCouriers() (*dto.AllCouriersResponse, error) {
	models, err := c.courierRepository.GetAllCourier()
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	var dtoResponse dto.AllCouriersResponse
	for _, data := range models {
		temp := dto.SingleCourierResponse{
			Id:   data.Id,
			Name: data.Name,
		}
		dtoResponse.Couriers = append(dtoResponse.Couriers, temp)
	}
	return &dtoResponse, nil
}
