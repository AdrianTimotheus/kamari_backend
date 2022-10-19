package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"strings"
	"time"
)

var (
	errOnlyRemoveWillComeVoucher = errors.New("only will come voucher is removable")
	errVoucherNotUpdatable       = errors.New("updatable voucher is only voucher that is ongoing or will come")
)

type MerchantVoucherService interface {
	GetVoucherDetail(voucherId uint64) (*dto.MerchantVoucherResponse, error)
	ValidateVoucherOwner(username string, voucherId uint64) error
	GetVouchers(merchantId uint64, keyword string, filters *helper.Filters) ([]dto.MerchantVoucherResponse, *helper.Metadata, error)
	CreateVoucher(payload *dto.MerchantVoucherRequest, username string) error
	UpdateVoucher(voucherId uint64, payload *dto.MerchantVoucherUpdate) error
	DeleteVoucher(voucherId uint64) error
	GetAllActiveVoucher(merchantId uint64) ([]dto.MerchantVoucherResponse, error)
}

type merchantVoucherService struct {
	merchantVoucherRepository repository.MerchantVoucherRepository
	merchantRepository        repository.MerchantRepository
}

type MVSConfig struct {
	MerchantVoucherRepository repository.MerchantVoucherRepository
	MerchantRepository        repository.MerchantRepository
}

func NewMerchantVoucherService(c MVSConfig) MerchantVoucherService {
	return merchantVoucherService{
		merchantVoucherRepository: c.MerchantVoucherRepository,
		merchantRepository:        c.MerchantRepository,
	}
}

func (m merchantVoucherService) GetVoucherDetail(voucherId uint64) (*dto.MerchantVoucherResponse, error) {
	v, err := m.merchantVoucherRepository.GetVoucherDetail(voucherId)
	if err != nil {
		return nil, err
	}

	voucher := dto.MerchantVoucherResponse{
		Id:             v.Id,
		Title:          v.Title,
		Code:           v.Code,
		DiscountAmount: v.DiscountAmount,
		MinimumPayment: v.MinimumPayment,
		StartDate:      v.StartDate,
		EndDate:        v.EndDate,
		StockLeft:      v.StockLeft,
		MerchantId:     v.MerchantId,
		PromoTypeId:    v.PromoTypeId,
		PromoType:      v.PromoType,
	}

	return &voucher, nil
}

func (m merchantVoucherService) CreateVoucher(payload *dto.MerchantVoucherRequest, username string) error {
	var (
		startDate time.Time
		endDate   time.Time
	)

	if payload.StartDate.Int64 > payload.EndDate.Int64 {
		return httperror.BadRequestError("end date cannot be smaller than start date", "INVALID_DATE")
	} else {
		startDate = time.Unix(payload.StartDate.Int64, 0)
		endDate = time.Unix(payload.EndDate.Int64, 0)
	}

	if strings.ToUpper(payload.Prefix) != strings.ToUpper(username[:4]) {
		return httperror.BadRequestError("invalid code prefix", "INVALID_PREFIX_CODE")
	}

	merchant, err := m.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		return err
	}

	voucher := model.MerchantVoucher{
		Title:          payload.Title,
		Code:           strings.ToUpper(payload.Prefix + payload.Code),
		PromoTypeId:    payload.PromoTypeId,
		DiscountAmount: payload.DiscountAmount,
		StartDate:      startDate,
		EndDate:        endDate,
		StockLeft:      payload.StockLeft,
		MerchantId:     merchant.Id,
		MinimumPayment: payload.MinimumPayment,
	}

	if exist, err := m.merchantVoucherRepository.IsVoucherCodeExist(voucher.Code, voucher.MerchantId); err != nil || exist {
		return httperror.BadRequestError("voucher already exist", "VOUCHER_EXIST")
	}

	if err = m.merchantVoucherRepository.CreateMerchantVoucher(&voucher); err != nil {
		return err
	}

	return nil
}

func (m merchantVoucherService) ValidateVoucherOwner(username string, voucherId uint64) error {
	merchant, err := m.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		return err
	}

	valid, err := m.merchantVoucherRepository.ValidVoucherOwner(merchant.Id, voucherId)
	if err != nil || !valid {
		return httperror.ForbiddenError()
	}

	return nil
}

func (m merchantVoucherService) DeleteVoucher(voucherId uint64) error {
	voucher := model.MerchantVoucher{Id: voucherId}
	removable, err := m.merchantVoucherRepository.IsRemovable(voucher)
	if err != nil || !removable {
		return errOnlyRemoveWillComeVoucher
	}
	if err = m.merchantVoucherRepository.DeleteVoucher(&voucher); err != nil {
		return err
	}

	return nil
}

func (m merchantVoucherService) UpdateVoucher(voucherId uint64, payload *dto.MerchantVoucherUpdate) error {
	var (
		startDate time.Time
		endDate   time.Time
	)

	if payload.StartDate.Valid {
		startDate = time.Unix(payload.StartDate.Int64, 0)
	}
	if payload.EndDate.Valid {
		endDate = time.Unix(payload.EndDate.Int64, 0)
	}

	voucher := model.MerchantVoucher{
		Id:             voucherId,
		Title:          payload.Title,
		PromoTypeId:    payload.PromoTypeId,
		DiscountAmount: payload.DiscountAmount,
		MinimumPayment: payload.MinimumPayment,
		StartDate:      startDate,
		EndDate:        endDate,
		StockLeft:      payload.StockLeft,
	}

	valid, err := m.merchantVoucherRepository.IsUpdatable(voucher)
	if err != nil || !valid {
		return errVoucherNotUpdatable
	}

	err = m.merchantVoucherRepository.UpdateVoucher(&voucher)
	if err != nil {
		return err
	}

	return nil
}

func (m merchantVoucherService) GetVouchers(merchantId uint64, keyword string, filters *helper.Filters) ([]dto.MerchantVoucherResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 10)
	filters.Page = helper.ReadInt(filters.Page, 1)

	vouchers, metadata, err := m.merchantVoucherRepository.GetMerchantVouchers(merchantId, keyword, filters)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil, httperror.BadRequestError(err.Error(), "NO_DATA")
		default:
			return nil, nil, err
		}
	}

	res := make([]dto.MerchantVoucherResponse, len(vouchers))
	for i, v := range vouchers {
		res[i] = dto.MerchantVoucherResponse{
			Id:             v.Id,
			Title:          v.Title,
			Code:           v.Code,
			DiscountAmount: v.DiscountAmount,
			MinimumPayment: v.MinimumPayment,
			StartDate:      v.StartDate,
			EndDate:        v.EndDate,
			StockLeft:      v.StockLeft,
			MerchantId:     v.MerchantId,
			PromoTypeId:    v.PromoTypeId,
			PromoType:      v.PromoType,
		}
	}

	return res, metadata, nil
}

func (m merchantVoucherService) GetAllActiveVoucher(merchantId uint64) ([]dto.MerchantVoucherResponse, error) {
	models, err := m.merchantVoucherRepository.GetAllActiveVoucher(merchantId)
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
			Id:             data.Id,
			Title:          data.Title,
			Code:           data.Code,
			StockLeft:      data.StockLeft,
			DiscountAmount: data.DiscountAmount,
			MinimumPayment: data.MinimumPayment,
			StartDate:      data.StartDate,
			EndDate:        data.EndDate,
			MerchantId:     data.MerchantId,
			PromoTypeId:    data.PromoTypeId,
			PromoType: model.PromoType{
				Id:   data.PromoType.Id,
				Type: data.PromoType.Type,
			},
		}
		dtoResponse = append(dtoResponse, temp)
	}
	return dtoResponse, nil

}
