package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/model"
)

type WalletRepository interface {
	GetMyWalletDetail(userID uint64) (*model.Wallet, error)
	GetMyWalletDetailByWalletID(WalletId uint64) (*model.Wallet, error)
	ChangeUserWalletPin(userId uint64, Password string) (bool, error)
	CountUserWallet(userId uint64) (int64, error)
	CreateWallet(model model.Wallet) (*model.Wallet, error)
	ValidatePin(userId uint64, HashedPIN string) bool
	TopUpWallet(transaction *model.WalletTransaction, wModel *model.Wallet) error
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

func (w walletRepository) ChangeUserWalletPin(userId uint64, Password string) (bool, error) {
	var wModel model.Wallet
	db := w.db.Model(&wModel).
		Where("user_id=?", userId).
		Update("pin", Password)
	if db.Error != nil {
		return false, db.Error
	}
	if db.RowsAffected < 1 {
		return false, errors.New("no data change")
	}
	return true, nil
}

func (w walletRepository) CountUserWallet(userId uint64) (int64, error) {
	var wModel model.Wallet
	var count int64
	if err := w.db.Model(wModel).
		Where("user_id=?", userId).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (w walletRepository) CreateWallet(inModel model.Wallet) (*model.Wallet, error) {
	err := w.db.Create(&inModel).First(&inModel).Error
	if err != nil {
		return nil, err
	}
	return &inModel, nil

}

func (w walletRepository) ValidatePin(userId uint64, HashedPIN string) bool {
	var numberOfRecord int64
	var walletModel model.Wallet
	w.db.Model(&walletModel).Where("user_id=?", userId).
		Where("pin=?", HashedPIN).
		Count(&numberOfRecord)
	fmt.Println(numberOfRecord)
	if numberOfRecord != 1 {
		return false
	}
	return true
}

func (w walletRepository) TopUpWallet(transaction *model.WalletTransaction, wModel *model.Wallet) error {
	err := w.db.Create(&transaction).Error
	if err != nil {
		return err
	}

	db := w.db.Model(&wModel).
		Where("user_id=?", wModel.UserId).
		Update("balance", gorm.Expr("balance + ?", transaction.Amount))
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		return errors.New("no data change")
	}
	return nil
}
