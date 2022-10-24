package service

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io/ioutil"
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/repository"
	"net/http"
	"strings"
)

type CourierService interface {
	GetAllCouriers() (*dto.AllCouriersResponse, error)
	PostCourierService(bodyReq string) ([]dto.CourierServiceRes, error)
	GetShippingCost(payload *dto.GetShippingCost) *dto.MerchantCourierDetailRes
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

func (c courierService) GetShippingCost(payload *dto.GetShippingCost) *dto.MerchantCourierDetailRes {
	jsonPayload, _ := json.Marshal(payload)
	resp, err := c.PostCourierService(string(jsonPayload))
	if err != nil {
		return nil
	}

	if len(resp) == 0 {
		return nil
	}

	data := resp[0]

	for _, cost := range data.Costs {
		if cost.Service == "REG" || cost.Service == "Pos Reguler" || cost.Service == "CTC" {
			c := dto.MerchantCourierDetailRes{Courier: payload.Courier}
			for _, cDetail := range cost.Cost {
				c.Etd = helper.ClearString(cDetail.Etd)
				c.Price = float64(cDetail.Value)
			}
			return &c
		}
	}

	return nil
}

func (c *courierService) PostCourierService(bodyReq string) ([]dto.CourierServiceRes, error) {
	sConfig := config.Config.CourierService
	req, _ := http.NewRequest("POST", "https://ship.spacexkode.my.id/services", strings.NewReader(bodyReq))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Api-Key", sConfig.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	res := make([]dto.CourierServiceRes, 0)
	json.Unmarshal(body, &res)
	return res, err
}
