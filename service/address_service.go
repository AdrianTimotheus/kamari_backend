package service

import (
	"kamari.com/dto"
	"kamari.com/model"
	"kamari.com/repository"
)

type AddressService interface {
	GetAddress(userId uint64) ([]model.Address, error)
	EditAddress(userId uint64, req dto.EditAddressReq) (*model.Address, error)
	AddNewAddress(userId uint64, req dto.AddNewAddressRequest) (*model.Address, error)
	DeleteAddress(userId uint64, id int) error
	GetDefaultAddress(userId uint64) (*model.Address, error)
	SetDefaultAddress(userId uint64, id int) (*model.Address, error)
}

type addressService struct {
	addressRepository repository.AddressRepository
}

type ASConfig struct {
	AddressRepository repository.AddressRepository
}

func NewAddressService(c ASConfig) AddressService {
	return addressService{
		addressRepository: c.AddressRepository,
	}
}

func (a addressService) GetAddress(userId uint64) ([]model.Address, error) {
	response, err := a.addressRepository.GetAddress(userId)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a addressService) EditAddress(userId uint64, req dto.EditAddressReq) (*model.Address, error) {
	address := model.Address{
		Id:            req.Id,
		Title:         req.Title,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		CityId:        req.CityId,
		DetailAddress: req.DetailAddress,
		IsDefault:     false,
	}
	response, err := a.addressRepository.EditAddress(userId, address)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a addressService) AddNewAddress(userId uint64, req dto.AddNewAddressRequest) (*model.Address, error) {
	address := model.Address{
		UserId:        userId,
		Title:         req.Title,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		CityId:        req.CityId,
		DetailAddress: req.DetailAddress,
		IsDefault:     false,
	}
	response, err := a.addressRepository.AddNewAddress(address)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a addressService) DeleteAddress(userId uint64, id int) error {
	err := a.addressRepository.DeleteAddress(userId, id)
	if err != nil {
		return err
	}
	return nil
}

func (a addressService) GetDefaultAddress(userId uint64) (*model.Address, error) {
	response, err := a.addressRepository.GetDefaultAddress(userId)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a addressService) SetDefaultAddress(userId uint64, id int) (*model.Address, error) {
	response, err := a.addressRepository.SetDefaultAddress(userId, id)
	if err != nil {
		return nil, err
	}
	return response, nil
}
