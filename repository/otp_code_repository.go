package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/model"
)

type OtpCodeRepository interface {
	ClearUserOTP(UserId uint64) (bool, error)
	CreateOTP(model model.OtpCode) (*model.OtpCode, error)
}

type otpCodeRepository struct {
	db *gorm.DB
}

type OCRConfig struct {
	DB *gorm.DB
}

func NewOtpCodeRepository(c OCRConfig) OtpCodeRepository {
	return otpCodeRepository{db: c.DB}
}

func (o otpCodeRepository) ClearUserOTP(UserId uint64) (bool, error) {
	fmt.Println("repos")
	err := o.db.Raw("DELETE FROM otp_codes WHERE user_id=?", UserId).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o otpCodeRepository) CreateOTP(model model.OtpCode) (*model.OtpCode, error) {
	err := o.db.Create(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, err
}
