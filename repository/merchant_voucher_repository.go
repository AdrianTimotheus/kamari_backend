package repository

import "gorm.io/gorm"

type MerchantVoucherRepository interface {
}

type merchantVoucherRepository struct {
	db *gorm.DB
}

type MVRConfig struct {
	DB *gorm.DB
}

func NewMerchantVoucherRepository(c MVRConfig) MerchantRepository {
	return merchantRepository{
		db: c.DB,
	}
}

//func (m merchantVoucherRepository) DeleteVoucher(voucherId uint64) error {
//	m.db.
//		Where("id = ?", voucherId).
//		Delete()
//}
