package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
	"time"
)

type TransactionRepository interface {
	CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, marketplaceVoucher *model.MarketplaceVoucher, merchantVoucher []model.MerchantVoucher, walletTransaction *model.WalletTransaction, promoItems []model.ProductPromoItem, productVariantModels []model.ProductVariant) (*model.Transaction, error)
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

func (t transactionRepository) CreateTransactionPayWithWallet(transaction *model.Transaction, wallet *model.Wallet, marketplaceVoucher *model.MarketplaceVoucher, merchantVoucher []model.MerchantVoucher, walletTransaction *model.WalletTransaction, productPromoItems []model.ProductPromoItem, productVariantModels []model.ProductVariant) (*model.Transaction, error) {
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

		for _, data := range productVariantModels {
			err := tx.Model(&data).Where("product_variants.id=?", data.ID).
				Update("stock", gorm.Expr("stock - ?", data.Stock)).Scan(&data).Error
			if err != nil {
				return httperror.InternalServerError("Internal Server Error")
			}
			if data.Stock < 0 {
				return httperror.BadRequestError("Insufficient Product Stock", "INSUFFICIENT_PRODUCT_STOCK")
			}
		}

		for _, data := range productPromoItems {
			var resp model.ProductPromoItem
			now := time.Now().Format("2006-01-02")
			err := tx.Model(&resp).
				Joins("ProductPromo").
				Where("product_id=?", data.ProductId).
				Where("start_date <= ?", now).
				Where("end_date >= ?", now).
				First(&resp).Error
			if err != nil {
				fmt.Println(err)
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return httperror.InternalServerError("Internal Server Error")
				}
			} else {
				fmt.Println("MAX: ", resp.MaximumItem)
				fmt.Println("taken: ", data.StockLeft)

				if resp.MaximumItem >= data.StockLeft {
					fmt.Println("MAX: ", resp.MaximumItem)
					fmt.Println("taken: ", data.StockLeft)
					err := tx.Model(&resp).
						Where("id=?", resp.Id).
						Update("stock_left", gorm.Expr("stock_left - ?", data.StockLeft)).
						First(&resp).Error
					if err != nil {
						return httperror.BadRequestError("Out of product promotion. Recalculate Total Payment", "BAD_REQUEST_OUT_OF_PRODUCT_PROMOTION")
					}
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
