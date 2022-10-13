package repository

import (
	"time"

	"gorm.io/gorm"
	"kamari.com/model"
)

type MerchantVoucherRepository interface {
	DeleteVoucher(voucher *model.MerchantVoucher) error
	UpdateVoucher(voucher *model.MerchantVoucher) error
	ValidVoucherOwner(merchantId, voucherId uint64) (bool, error)
	IsRemovable(voucher model.MerchantVoucher) (bool, error)
	IsUpdatable(voucher model.MerchantVoucher) (bool, error)
}

type merchantVoucherRepository struct {
	db *gorm.DB
}

type MVRConfig struct {
	DB *gorm.DB
}

func NewMerchantVoucherRepository(c MVRConfig) MerchantVoucherRepository {
	return merchantVoucherRepository{
		db: c.DB,
	}
}

func (m merchantVoucherRepository) ValidVoucherOwner(merchantId, voucherId uint64) (bool, error) {
	var voucher model.MerchantVoucher
	res := m.db.
		Where("id = ?", voucherId).
		Where("merchant_id = ?", merchantId).
		Limit(1).
		Find(&voucher)

	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (m merchantVoucherRepository) IsRemovable(voucher model.MerchantVoucher) (bool, error) {
	now := time.Now().Format("2006-01-02")
	res := m.db.
		Where("id = ?", voucher.Id).
		Where("start_date > ?", now).
		Limit(1).
		Find(&voucher)
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (m merchantVoucherRepository) DeleteVoucher(voucher *model.MerchantVoucher) error {
	return m.db.Delete(&voucher).Error
}

func (m merchantVoucherRepository) UpdateVoucher(voucher *model.MerchantVoucher) error {
	return m.db.Updates(&voucher).Error
}

func (m merchantVoucherRepository) IsUpdatable(voucher model.MerchantVoucher) (bool, error) {
	now := time.Now().Format("2006-01-02")
	res := m.db.
		Where("id = ?", voucher.Id).
		Where("start_date >= ?", now).
		Or("now() < end_date").
		Limit(1).
		Find(&voucher)
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}
