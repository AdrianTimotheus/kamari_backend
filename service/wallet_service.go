package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type WalletService interface {
	GetUserWalletDetail(userModel model.User) (*dto.DetailWalletResponse, error)
}

type walletService struct {
	walletRepository repository.WalletRepository
}

type WSConfig struct {
	WalletRepository repository.WalletRepository
}

func NewWalletService(c WSConfig) WalletService {
	return walletService{walletRepository: c.WalletRepository}
}

func (w walletService) GetUserWalletDetail(userModel model.User) (*dto.DetailWalletResponse, error) {
	wModel, err := w.walletRepository.GetMyWalletDetail(userModel.Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	dtoResponse := dto.DetailWalletResponse{
		Id:      wModel.Id,
		Balance: wModel.Balance,
	}
	return &dtoResponse, nil
}
