package repository

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/model"
)

type SealabsAccountRepository interface {
	CreateAccount(account *model.SealabsPayAccount) (uint64, error)
	ChangeDefaultTransaction(userId, accountId uint64) error
	HasSameCard(sealabsAccount *model.SealabsPayAccount) bool
	GetAccounts(userId uint64) ([]model.SealabsPayAccount, error)
	IsValidOwner(userId, accountId uint64) (bool, error)
	RemoveAccountTransaction(account *model.SealabsPayAccount) error
}

type sealabsAccountRepository struct {
	db *gorm.DB
}

type SARConfig struct {
	DB *gorm.DB
}

func NewSealabsAccountRepository(c SARConfig) SealabsAccountRepository {
	return sealabsAccountRepository{
		db: c.DB,
	}
}

func (s sealabsAccountRepository) CreateAccount(account *model.SealabsPayAccount) (uint64, error) {
	err := s.db.Create(&account).Error
	if err != nil {
		return 0, err
	}

	return account.Id, nil
}

func (s sealabsAccountRepository) ChangeDefaultTransaction(userId, accountId uint64) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.disableDefault(userId, tx)
		if err != nil {
			return err
		}

		err = s.setAccountToDefault(accountId, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s sealabsAccountRepository) disableDefault(userId uint64, tx *gorm.DB) error {
	err := tx.Model(model.SealabsPayAccount{}).
		Where("user_id", userId).
		Where("is_default = ?", true).
		Update("is_default", false).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (sealabsAccountRepository) setAccountToDefault(accountId uint64, tx *gorm.DB) error {
	return tx.Model(model.SealabsPayAccount{}).
		Where("id = ?", accountId).
		Update("is_default", true).
		Error
}

func (s sealabsAccountRepository) HasSameCard(sealabsAccount *model.SealabsPayAccount) bool {
	var acc model.SealabsPayAccount
	err := s.db.
		Where("user_id = ?", sealabsAccount.UserId).
		Where("title = ?", sealabsAccount.Title).
		Or("user_id = ?", sealabsAccount.UserId).
		Where("card_number = ?", sealabsAccount.CardNumber).
		First(&acc).Error

	return err == nil
}

func (s sealabsAccountRepository) GetAccounts(userId uint64) ([]model.SealabsPayAccount, error) {
	var accounts []model.SealabsPayAccount

	err := s.db.
		Where("user_id = ?", userId).
		Order("is_default DESC, id ASC").
		Find(&accounts).
		Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s sealabsAccountRepository) IsValidOwner(userId, accountId uint64) (bool, error) {
	res := s.db.
		Where("user_id = ?", userId).
		Where("id = ?", accountId).
		Limit(1).
		Find(&model.SealabsPayAccount{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (s sealabsAccountRepository) RemoveAccountTransaction(account *model.SealabsPayAccount) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.deleteSealabsAccount(tx, account)
		if err != nil {
			return err
		}

		var randomAccount model.SealabsPayAccount
		err = s.db.First(&randomAccount).Error
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return nil
			default:
				return err
			}
		}

		if err := s.setAccountToDefault(randomAccount.Id, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return s.db.Delete(&account).Error
}

func (sealabsAccountRepository) deleteSealabsAccount(tx *gorm.DB, account *model.SealabsPayAccount) error {
	return tx.Delete(&account).Error
}
