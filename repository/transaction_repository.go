package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type TransactionRepository interface {
	CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, walletTransaction *model.WalletTransaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

type TRConfig struct {
	DB *gorm.DB
}

func NewTransactionRepository(c TRConfig) TransactionRepository {
	return &transactionRepository{db: c.DB}
}

func (t transactionRepository) CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, walletTransaction *model.WalletTransaction) error {
	err := t.db.Transaction(func(tx *gorm.DB) error {
		var updatedTrcModel model.Transaction

		if err := tx.Create(&transaction).First(&updatedTrcModel).Error; err != nil {
			return err
		}

		if err := tx.Model(&wallet).Update("balance", gorm.Expr("balance-?", walletTransaction.Amount)).Error; err != nil {
			return err
		}

		if err := tx.Exec("INSERT INTO wallet_transactions (src_wallet_id, dst_wallet_id, transaction_type_id, amount, description) values (?,?,?,?,?)"+
			"", wallet.Id, nil, walletTransaction.TransactionTypeId, walletTransaction.Amount, walletTransaction.Description).Error; err != nil {
			return err
		}

		walletPayModel := model.WalletPayment{
			WalletId:      wallet.Id,
			TransactionId: updatedTrcModel.Id,
		}

		if err := tx.Model(&walletPayModel).Create(&walletPayModel).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
