package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type WalletRepository interface {
	GetMyWalletDetail(userID uint64) (*model.Wallet, error)
	GetMyWalletDetailByWalletID(WalletId uint64) (*model.Wallet, error)
}

type walletRepository struct {
	db *gorm.DB
}

type WConfig struct {
	DB *gorm.DB
}

func NewWalletRepository(c WConfig) WalletRepository {
	return walletRepository{db: c.DB}
}

func (w walletRepository) GetMyWalletDetail(userID uint64) (*model.Wallet, error) {
	var wModel model.Wallet
	err := w.db.
		Joins("User").
		Where("user_id=?", userID).
		First(&wModel).Error
	if err != nil {
		return nil, err
	}
	return &wModel, nil
}

func (w walletRepository) GetMyWalletDetailByWalletID(WalletId uint64) (*model.Wallet, error) {
	var wModel model.Wallet
	err := w.db.
		Joins("User").
		Where("wallets.id=?", WalletId).
		First(&wModel).Error
	if err != nil {
		return nil, err
	}
	return &wModel, nil
}
