package service

import (
	"errors"

	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type MerchantVoucherService interface {
	ValidateVoucherOwner(username string, voucherId uint64) error
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
	removable, err := m.merchantVoucherRepository.IsRemovable(&voucher)
	if err != nil || !removable {
		return errors.New("only remove will come voucher")
	}
	if err = m.merchantVoucherRepository.DeleteVoucher(&voucher); err != nil {
		return err
	}

	return nil
}
