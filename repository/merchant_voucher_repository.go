package repository

import (
	"fmt"
	"kamari.com/helper"
	"time"

	"gorm.io/gorm"
	"kamari.com/model"
)

type MerchantVoucherRepository interface {
	GetVoucherDetail(voucherId uint64) (*model.MerchantVoucher, error)
	DeleteVoucher(voucher *model.MerchantVoucher) error
	UpdateVoucher(voucher *model.MerchantVoucher) error
	CreateMerchantVoucher(voucher *model.MerchantVoucher) error
	GetMerchantVouchers(merchantId uint64, keyword string, filters *helper.Filters) ([]model.MerchantVoucher, *helper.Metadata, error)
	ValidVoucherOwner(merchantId, voucherId uint64) (bool, error)
	IsRemovable(voucher model.MerchantVoucher) (bool, error)
	IsUpdatable(voucher model.MerchantVoucher) (bool, error)
	IsVoucherCodeExist(voucherCode string, merchantId uint64) (bool, error)
	GetAllActiveVoucher(merchantId uint64) ([]model.MerchantVoucher, error)
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

func (m merchantVoucherRepository) GetVoucherDetail(voucherId uint64) (*model.MerchantVoucher, error) {
	var voucher model.MerchantVoucher
	err := m.db.
		Preload("PromoType").
		Where("id = ?", voucherId).
		First(&voucher).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
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
		Where("start_date > ?", now).
		Or("start_date >= ?", now).
		Where("now() < end_date").
		Limit(1).
		Find(&voucher)
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (m merchantVoucherRepository) GetMerchantVouchers(merchantId uint64, keyword string, filters *helper.Filters) ([]model.MerchantVoucher, *helper.Metadata, error) {
	var (
		vouchers []model.MerchantVoucher
		count    int64
	)

	k := fmt.Sprintf("%s%%", keyword)
	res := m.db.
		Preload("PromoType").
		Where("merchant_id = ?", merchantId).
		Where("title ILIKE ? OR code ILIKE ?", k, k)

	if filters.VoucherStatus != nil {
		now := time.Now().Format("2006-01-02")
		switch *filters.VoucherStatus {
		case "ongoing":
			res = res.
				Where("start_date <= ?", now).
				Where("end_date >= ?", now)
		case "willCome":
			res = res.
				Where("start_date > ?", now)
		case "hasEnded":
			res = res.
				Where("end_date < ?", now)
		}
	}

	if err := res.Find(&vouchers).Count(&count).Error; err != nil {
		return nil, nil, err
	}

	res = res.
		Order("updated_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&vouchers)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)

	return vouchers, &metadata, nil
}

func (m merchantVoucherRepository) CreateMerchantVoucher(voucher *model.MerchantVoucher) error {
	return m.db.Create(&voucher).Error
}

func (m merchantVoucherRepository) IsVoucherCodeExist(voucherCode string, merchantId uint64) (bool, error) {
	var voucher model.MerchantVoucher
	res := m.db.
		Where("code = ?", voucherCode).
		Where("merchant_id = ?", merchantId).
		Limit(1).
		Find(&voucher)

	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (m merchantVoucherRepository) GetAllActiveVoucher(merchantId uint64) ([]model.MerchantVoucher, error) {
	var models []model.MerchantVoucher
	now := time.Now().Format("2006-01-02")
	err := m.db.
		Where("merchant_id=?", merchantId).
		Where("start_date <= ?", now).
		Where("end_date >= ?", now).
		Where("stock_left>0").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
