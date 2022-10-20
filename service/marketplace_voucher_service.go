package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type MarketplaceVoucherService interface {
	GetAllActiveMarketplaceVoucher() ([]dto.MerchantVoucherResponse, error)
}

type marketplaceVoucherService struct {
	marketplaceVoucherRepository repository.MarketplaceVoucherRepository
}

type MpVSConfig struct {
	MarketplaceVoucherRepository repository.MarketplaceVoucherRepository
}

func NewMarketplaceVoucherService(config MpVSConfig) MarketplaceVoucherService {
	return marketplaceVoucherService{marketplaceVoucherRepository: config.MarketplaceVoucherRepository}
}

func (m marketplaceVoucherService) GetAllActiveMarketplaceVoucher() ([]dto.MerchantVoucherResponse, error) {
	models, err := m.marketplaceVoucherRepository.GetAllActiveVoucher()

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	var dtoResponse []dto.MerchantVoucherResponse
	for _, data := range models {
		temp := dto.MerchantVoucherResponse{
			Id:             uint64(data.ID),
			Title:          data.Title,
			Code:           data.Code,
			StockLeft:      uint(data.StockLeft),
			DiscountAmount: data.DiscountAmount,
			StartDate:      data.StartDate,
			EndDate:        data.EndDate,
			PromoTypeId:    data.PromoTypeID,
			PromoType: model.PromoType{
				Id:   data.PromoType.Id,
				Type: data.PromoType.Type,
			},
		}
		dtoResponse = append(dtoResponse, temp)
	}
	return dtoResponse, nil
}
