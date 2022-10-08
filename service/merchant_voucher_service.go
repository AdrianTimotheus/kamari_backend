package service

import "kamari.com/repository"

type MerchantVoucherService interface {
}

type merchantVoucherService struct {
	merchantVoucherRepository repository.MerchantVoucherRepository
}

type MVSConfig struct {
	MerchantVoucherRepository repository.MerchantVoucherRepository
}

func NewMerchantVoucherService(c MVSConfig) MerchantVoucherService {
	return merchantVoucherService{
		merchantVoucherRepository: c.MerchantVoucherRepository,
	}
}

func (m merchantVoucherService) DeleteVoucher() {

}
