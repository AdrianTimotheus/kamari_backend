package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type WalletService interface {
	GetUserWalletDetail(userModel model.User) (*dto.DetailWalletResponse, error)
	GetWalletTransactionHistories(userModel model.User, payload dto.WalletTransactionRequest) (*dto.WalletTransactionResponse, error)
	ChangeUserWalletPin(userId uint64, Password string) (bool, error)
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
	tempWalletData := make(map[uint64]*model.Wallet)
	for _, data := range wtModels {
		temp := dto.TransactionRecord{
			Id: data.Id,
		}
		if data.SrcWalletId != 0 {
			_, ok := tempWalletData[data.SrcWalletId]
			if !ok {
				tempWModel, _ := w.walletRepository.GetMyWalletDetailByWalletID(data.SrcWalletId)
				tempWalletData[tempWModel.Id] = tempWModel
			}
			temp.Source.WalletId = tempWalletData[data.SrcWalletId].Id
			temp.Source.Name = tempWalletData[data.SrcWalletId].User.Name
		}
		if data.DstWalletId != 0 {
			_, ok := tempWalletData[data.DstWalletId]
			if !ok {
				tempWModel, _ := w.walletRepository.GetMyWalletDetailByWalletID(data.DstWalletId)
				tempWalletData[tempWModel.Id] = tempWModel
			}
			temp.Destination.WalletId = tempWalletData[data.DstWalletId].Id
			temp.Destination.Name = tempWalletData[data.DstWalletId].User.Name
		}

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

func (w walletService) ChangeUserWalletPin(userId uint64, Password string) (bool, error) {
	hashedPassword := fmt.Sprintf("%X", sha256.Sum256([]byte(Password)))
	status, err := w.walletRepository.ChangeUserWalletPin(userId, hashedPassword)
	if err != nil {
		return false, err
	}
	return status, nil
}
