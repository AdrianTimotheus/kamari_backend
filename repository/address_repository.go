package repository

import (
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
)

type AddressRepository interface {
	GetAddress(userId uint64) ([]model.Address, error)
	EditAddress(userId uint64, address model.Address) (*model.Address, error)
	AddNewAddress(address model.Address) (*model.Address, error)
	DeleteAddress(userId uint64, id int) error
	GetDefaultAddress(userId uint64) (*model.Address, error)
	SetDefaultAddress(userId uint64, id int) (*model.Address, error)
}

type addressRepository struct {
	db *gorm.DB
}

type ARConfig struct {
	DB *gorm.DB
}

func NewAddressRepository(c ARConfig) AddressRepository {
	return addressRepository{db: c.DB}
}

func (a addressRepository) GetAddress(userId uint64) ([]model.Address, error) {
	var addresses []model.Address
	err := a.db.Table("user_addresses").Where("user_id = ?", userId).Find(&addresses).Error
	if err != nil {
		return nil, httperror.InternalServerError("Failed to get user address")
	}
	return addresses, nil
}

func (a addressRepository) EditAddress(userId uint64, address model.Address) (*model.Address, error) {
	var userAddress model.Address
	var exists bool
	var err error

	err = a.db.Table("user_addresses").
		Select("count(*) > 0").
		Where("id = ?", address.Id).
		Where("user_id = ?", userId).
		Find(&exists).
		Error

	if exists {
		err = a.db.Table("user_addresses").
			Where("user_id = ?", userId).
			Where("id = ?", address.Id).
			Updates(&address).
			First(&userAddress).
			Error

		if err != nil {
			return nil, httperror.InternalServerError("Failed to update user address")
		}

		return &userAddress, nil
	}

	return nil, httperror.BadRequestError("Address not found", "ADDRESS_NOT_FOUND")
}

func (a addressRepository) AddNewAddress(address model.Address) (*model.Address, error) {
	var userAddress model.Address
	var exists bool

	//set default address if first address
	err := a.db.Table("user_addresses").
		Select("count(*) > 0").
		Where("user_id = ?", address.UserId).
		Find(&exists).Error

	if !exists {
		address.IsDefault = true
	} else {
		address.IsDefault = false
	}

	err = a.db.Table("user_addresses").Create(&address).Scan(&userAddress).Error

	if err != nil {
		return nil, httperror.InternalServerError("Failed to add new user address")
	}

	return &userAddress, nil
}

func (a addressRepository) DeleteAddress(userId uint64, id int) error {
	var userAddress model.Address
	var exists bool
	var err error

	err = a.db.Table("user_addresses").
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Find(&exists).
		Error

	if exists {
		err = a.db.Table("user_addresses").
			Where("user_id = ?", userId).
			Where("id = ?", id).
			Delete(&userAddress).
			Error

		if err != nil {
			return httperror.InternalServerError("Failed to delete user address")
		}

		return nil
	}

	return httperror.BadRequestError("Address not found", "ADDRESS_NOT_FOUND")
}

func (a addressRepository) GetDefaultAddress(userId uint64) (*model.Address, error) {
	var address model.Address

	err := a.db.Table("user_addresses").
		Where("user_id = ?", userId).
		Where("is_default = ?", true).
		First(&address).
		Error

	if err != nil {
		return nil, httperror.InternalServerError("Default address has not been set yet")
	}

	return &address, nil
}

func (a addressRepository) SetDefaultAddress(userId uint64, id int) (*model.Address, error) {
	var address model.Address
	var err error
	var exists bool

	err = a.db.Table("user_addresses").
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Find(&exists).
		Error

	if !exists {
		return nil, httperror.BadRequestError("Address not found", "ADDRESS_NOT_FOUND")
	}

	err = a.db.Table("user_addresses").
		Select("count(*) > 0").
		Where("is_default = ?", true).
		Where("user_id = ?", userId).
		Find(&exists).
		Error

	if exists {
		err = a.db.Table("user_addresses").
			Where("is_default = ?", true).
			Where("user_id = ?", userId).
			Update("is_default", false).
			Error

		if err != nil {
			return nil, httperror.InternalServerError("Failed to update default address")
		}
	}

	err = a.db.Table("user_addresses").
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Update("is_default", true).
		Scan(&address).
		Error

	if err != nil {
		return nil, httperror.InternalServerError("Failed to update default address")
	}

	return &address, nil
}
