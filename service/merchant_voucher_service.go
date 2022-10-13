package service

import (
	"errors"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"time"
)

var (
	errOnlyRemoveWillComeVoucher = errors.New("only will come voucher is removable")
	errVoucherNotUpdatable       = errors.New("updatable voucher is only voucher that is ongoing or will come")
)

type MerchantVoucherService interface {
	ValidateVoucherOwner(username string, voucherId uint64) error
	UpdateVoucher(voucherId uint64, payload *dto.MerchantVoucherUpdate) error
	DeleteVoucher(voucherId uint64) error
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
