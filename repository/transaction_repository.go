package repository

import (
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
)

type TransactionRepository interface {
	CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, marketplaceVoucher *model.MarketplaceVoucher, merchantVoucher []model.MerchantVoucher, walletTransaction *model.WalletTransaction) (*model.Transaction, error)
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

func (t transactionRepository) CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, marketplaceVoucher *model.MarketplaceVoucher, merchantVoucher []model.MerchantVoucher, walletTransaction *model.WalletTransaction) (*model.Transaction, error) {
	err := t.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&transaction).Error; err != nil {
			return httperror.InternalServerError("Internal Server Error")
		}

		if err := tx.Model(&wallet).Update("balance", gorm.Expr("balance-?", walletTransaction.Amount)).Error; err != nil {
			return httperror.InternalServerError("Internal Server Error")
		}

		if marketplaceVoucher.ID != 0 {
			if err := tx.Model(&marketplaceVoucher).Update("stock_left", gorm.Expr("stock_left-1")).First(&marketplaceVoucher).Error; err != nil {
				return httperror.BadRequestError("Out of Marketplace Voucher", "OUT_OF_MARKETPLACE_VOUCHER")
				//return err
			}
			if marketplaceVoucher.StockLeft < 0 {
				return httperror.BadRequestError("Out of Marketplace Voucher", "OUT_OF_MARKETPLACE_VOUCHER")
			}
		}

		for _, data := range merchantVoucher {
			if data.Id != 0 {
				if err := tx.Model(&data).Update("stock_left", gorm.Expr("stock_left-1")).First(&data).Error; err != nil {
					return httperror.BadRequestError("Out of Merchant Voucher", "OUT_OF_MERCHANT_VOUCHER")
					//return err
				}
				if data.StockLeft < 0 {
					return httperror.BadRequestError("Out of Merchant Voucher", "OUT_OF_MERCHANT_VOUCHER")
				}
			}
		}

		if err := tx.Model(&wallet).Update("balance", gorm.Expr("balance-?", walletTransaction.Amount)).Error; err != nil {
			return httperror.InternalServerError("Internal Server Error")
		}

		if err := tx.Exec("INSERT INTO wallet_transactions (src_wallet_id, dst_wallet_id, transaction_type_id, amount, description) values (?,?,?,?,?)"+
			"", wallet.Id, nil, walletTransaction.TransactionTypeId, walletTransaction.Amount, walletTransaction.Description).Error; err != nil {
			return httperror.InternalServerError("Internal Server Error")
		}

		walletPayModel := model.WalletPayment{
			WalletId:      wallet.Id,
			TransactionId: transaction.Id,
		}

		if err := tx.Model(&walletPayModel).Create(&walletPayModel).Error; err != nil {
			return httperror.InternalServerError("Internal Server Error")
		}
		return nil
	})

	return transaction, err
}
