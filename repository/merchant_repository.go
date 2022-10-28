package repository

import (
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
)

type MerchantRepository interface {
	GetMerchantDetail(username string) (model.Merchant, error)
	GetMerchantProfile(username string) (model.Merchant, error)
	UpdateMerchantProfile(merchant *model.Merchant) error
	GetMerchantDetailById(id uint64) (model.Merchant, error)
	RegisterMerchant(userId uint64, req *dto.MerchantRegisterRequest) error
	GetMerchantIncome(merchant model.Merchant) (uint64, error)
	WithdrawIncome(user model.User, merchant model.Merchant) (uint64, error)
}

type merchantRepository struct {
	db *gorm.DB
}

type MRConfig struct {
	DB *gorm.DB
}

func NewMerchantRepository(c MRConfig) MerchantRepository {
	return merchantRepository{
		db: c.DB,
	}
}

func (m merchantRepository) GetMerchantDetail(username string) (model.Merchant, error) {
	var merchant model.Merchant
	err := m.db.
		Joins("User").
		Joins("UserAddress").
		Where(`"User"."username" = ?`, username).
		First(&merchant).Error
	if err != nil {
		return merchant, err
	}

	return merchant, nil
}

func (m merchantRepository) GetMerchantProfile(username string) (model.Merchant, error) {
	var merchant model.Merchant
	err := m.db.
		Joins("User").
		Preload("UserAddress.City.Province").
		Where(`"User"."username" = ?`, username).
		First(&merchant).Error

	if err != nil {
		return merchant, err
	}

	return merchant, nil
}

func (m merchantRepository) UpdateMerchantProfile(merchant *model.Merchant) error {
	return m.db.
		Updates(&merchant).Error
}

func (m merchantRepository) GetMerchantDetailById(id uint64) (model.Merchant, error) {
	var merchant model.Merchant
	err := m.db.
		Joins("User").
		Joins("UserAddress").
		Where(`"merchants"."id" = ?`, id).
		First(&merchant).Error

	if err != nil {
		return merchant, err
	}

	return merchant, nil
}

func (m merchantRepository) RegisterMerchant(userId uint64, req *dto.MerchantRegisterRequest) error {
	var merchant model.Merchant
	var exists bool

	err := m.db.Model(&merchant).Select("count(*) > 0").Where("user_id = ?", userId).Find(&exists).Error
	if exists {
		return httperror.BadRequestError("User has already registered as a merchant", "ALREADY_MERCHANT_ERROR")
	}

	err = m.db.Model(&merchant).Select("count(*) > 0").Where("name = ?", req.Name).Find(&exists).Error
	if exists {
		return httperror.BadRequestError("Merchant name is already taken", "MERCHANT_NAME_TAKEN")
	}

	var userAddress model.UserAddress
	err = m.db.First(&userAddress).Where("user_id = ?", userId).Error
	if err != nil {
		return httperror.InternalServerError("Address not found")
	}

	newMerchant := model.Merchant{
		Name:          req.Name,
		Description:   req.Description,
		UserId:        userId,
		UserAddressId: userAddress.Id,
	}

	err = m.db.Table("users").Where("id = ?", userId).Update("role_id", model.RoleSeller).Error
	if err != nil {
		//return httperror.InternalServerError("Error updating user role")
		return err
	}

	err = m.db.Create(&newMerchant).Error
	if err != nil {
		return httperror.InternalServerError("Error creating merchant")
	}

	return nil
}

func (m merchantRepository) GetMerchantIncome(merchant model.Merchant) (uint64, error) {
	order := model.Order{MerchantId: merchant.Id}
	var income uint64
	err := m.db.Model(&order).Select("coalesce(SUM(merchant_income),0)").
		Where("merchant_id=?", merchant.Id).
		Where("order_status_id=?", 6).
		Where("withdraw_status=false").
		Scan(&income).Error
	if err != nil {
		return 0, err
	}
	return income, nil
}

func (m merchantRepository) WithdrawIncome(user model.User, merchant model.Merchant) (uint64, error) {
	order := model.Order{MerchantId: merchant.Id}

	var income uint64
	err := m.db.Model(&order).Select("SUM(merchant_income) ").
		Where("merchant_id=?", merchant.Id).
		Where("order_status_id=?", 6).
		Where("withdraw_status=false").
		Scan(&income).Error

	err = m.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&order).
			Where("merchant_id=?", merchant.Id).
			Where("order_status_id=?", 6).
			Update("withdraw_status", true).Error
		if err != nil {
			return err
		}

		var wallet model.Wallet
		err = tx.Model(&wallet).
			Where("user_id=?", user.Id).
			Update("balance", gorm.Expr("balance+?", income)).First(&wallet).
			Error
		if err != nil {
			return err
		}

		if income != 0 {
			walletTransaction := model.WalletTransaction{
				DstWalletId:       wallet.Id,
				Amount:            income,
				Description:       "Merchant income withdraw",
				TransactionTypeId: 2,
			}
			err = tx.Exec("INSERT INTO wallet_transactions (src_wallet_id, dst_wallet_id, transaction_type_id, amount, description) values (?,?,?,?,?)"+
				"", nil, wallet.Id, walletTransaction.TransactionTypeId, walletTransaction.Amount, walletTransaction.Description).Error
			//err = tx.Create(&walletTransaction).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return income, nil
}
