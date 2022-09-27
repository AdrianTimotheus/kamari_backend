package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type WalletTransactionRepository interface {
	GetUserWalletTransaction(walletId uint64, dataOffset int) ([]model.WalletTransaction, error)
	GetUserWalletTransactionByType(walletId uint64, dataOffset int, filter string) ([]model.WalletTransaction, error)
	CountWalletTransactionRecords(walletId uint64) (int64, error)
}

type walletTransaction struct {
	db *gorm.DB
}

type WTConfig struct {
	DB *gorm.DB
}

func NewWalletTransaction(c WTConfig) WalletTransactionRepository {
	return walletTransaction{db: c.DB}
}

func (w walletTransaction) GetUserWalletTransaction(walletId uint64, dataOffset int) ([]model.WalletTransaction, error) {
	var models []model.WalletTransaction
	err := w.db.
		Joins("TransactionType").
		Where("src_wallet_id = ? OR dst_wallet_id = ?", walletId, walletId).
		Order("id DESC").
		Offset(dataOffset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (w walletTransaction) GetUserWalletTransactionByType(walletId uint64, dataOffset int, filter string) ([]model.WalletTransaction, error) {
	var models []model.WalletTransaction
	err := w.db.
		Joins("TransactionType").
		Where("src_wallet_id = ? OR dst_wallet_id = ?", walletId, walletId).
		Where("transaction_type = ?", filter).
		Order("id DESC").
		Offset(dataOffset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (w walletTransaction) CountWalletTransactionRecords(walletId uint64) (int64, error) {
	var numOfRecords int64
	err := w.db.Table("wallet_transactions").
		Where("src_wallet_id = ? OR dst_wallet_id = ?", walletId, walletId).
		Count(&numOfRecords).Error
	if err != nil {
		return 0, err
	}
	return numOfRecords, nil

}
