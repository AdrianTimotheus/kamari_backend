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
	GetWalletTransactionHistories(userModel model.User, payload dto.WalletTransactionRequest) (*dto.WalletTransactionResponse, error)
}

type walletService struct {
	walletRepository            repository.WalletRepository
	walletTransactionRepository repository.WalletTransactionRepository
}

type WSConfig struct {
	WalletRepository            repository.WalletRepository
	WalletTransactionRepository repository.WalletTransactionRepository
}

func NewWalletService(c WSConfig) WalletService {
	return walletService{
		walletRepository:            c.WalletRepository,
		walletTransactionRepository: c.WalletTransactionRepository,
	}
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

func (w walletService) GetWalletTransactionHistories(userModel model.User, payload dto.WalletTransactionRequest) (*dto.WalletTransactionResponse, error) {
	wModel, err := w.walletRepository.GetMyWalletDetail(userModel.Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	recordsCount, err := w.walletTransactionRepository.CountWalletTransactionRecords(wModel.Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	dtoResponse := dto.WalletTransactionResponse{
		UserId:   userModel.Id,
		WalletID: wModel.Id,
		Page:     payload.Page,
	}
	if recordsCount%10 == 0 {
		dtoResponse.TotalPage = int(recordsCount / 10)
	} else {
		dtoResponse.TotalPage = int((recordsCount / 10) + 1)
	}
	DataOffset := (payload.Page * 10) - 10

	wtModels, err := w.walletTransactionRepository.GetUserWalletTransaction(wModel.Id, DataOffset)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	for _, data := range wtModels {
		var temp dto.TransactionRecord
		temp.Id = data.Id
		temp.SourceWallet = data.SrcWalletId
		temp.DestinationWallet = data.DstWalletId
		temp.Amount = data.Amount
		temp.Description = data.Description
		if data.SrcWalletId == wModel.Id {
			temp.TransactionType = data.TransactionType.Type
		} else {
			switch data.TransactionType.Type {
			case "DEBIT":
				temp.TransactionType = "CREDIT"
			case "CREDIT":
				temp.TransactionType = "DEBIT"
			}
		}
		temp.DateOfTime = data.CreatedAt.Format("01/02/2006")
		dtoResponse.Records = append(dtoResponse.Records, temp)
	}
	return &dtoResponse, nil
}
