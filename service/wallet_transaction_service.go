package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type WalletTransactionService interface {
	GetDetailTransaction(Id uint64, UserId uint64) (*dto.WalletTransactionDetailResponse, error)
	GetMerchantWithdrawHistories(user model.User) ([]dto.WithdrawHistoriesResponse, error)
}

type walletTransactionService struct {
	walletTransactionRepository repository.WalletTransactionRepository
	walletRepository            repository.WalletRepository
}

type WTSConfig struct {
	WalletTransactionRepository repository.WalletTransactionRepository
	WalletRepository            repository.WalletRepository
}

func NewWalletTransactionService(c WTSConfig) WalletTransactionService {
	return walletTransactionService{
		walletTransactionRepository: c.WalletTransactionRepository,
		walletRepository:            c.WalletRepository,
	}
}

func (w walletTransactionService) GetDetailTransaction(Id uint64, UserId uint64) (*dto.WalletTransactionDetailResponse, error) {
	wModel, err := w.walletRepository.GetMyWalletDetail(UserId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	wTModel, err := w.walletTransactionRepository.GetWalletTransactionByID(Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	if wModel.Id != wTModel.SrcWalletId && wModel.Id != wTModel.DstWalletId {
		return nil, httperror.ForbiddenError()
	}

	dtoResponse := dto.WalletTransactionDetailResponse{
		Amount:      wTModel.Amount,
		Description: wTModel.Description,
		DateOfTime:  wTModel.CreatedAt.Format("01/02/2006"),
	}

	srcWModel, err := w.walletRepository.GetMyWalletDetailByWalletID(wTModel.SrcWalletId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	dstWModel, err := w.walletRepository.GetMyWalletDetailByWalletID(wTModel.DstWalletId)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}

	if wTModel.SrcWalletId == wModel.Id {
		dtoResponse.TransactionType = wTModel.TransactionType.Type
	} else {
		switch wTModel.TransactionType.Type {
		case "DEBIT":
			dtoResponse.TransactionType = "CREDIT"
		case "CREDIT":
			dtoResponse.TransactionType = "DEBIT"
		}
	}
	dtoResponse.Source.WalletId = srcWModel.Id
	dtoResponse.Source.Name = srcWModel.User.Name
	dtoResponse.Destination.WalletId = dstWModel.Id
	dtoResponse.Destination.Name = dstWModel.User.Name
	return &dtoResponse, err
}

func (w walletTransactionService) GetMerchantWithdrawHistories(user model.User) ([]dto.WithdrawHistoriesResponse, error) {
	wallet, err := w.walletRepository.GetMyWalletDetail(user.Id)
	if err != nil {
		return nil, err
	}
	respModel, err := w.walletTransactionRepository.GetMerchantWithdrawHistories(wallet.Id)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	var dtoResponse []dto.WithdrawHistoriesResponse
	for _, data := range respModel {
		temp := dto.WithdrawHistoriesResponse{
			Amount: data.Amount,
			Time:   data.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		dtoResponse = append(dtoResponse, temp)
	}
	return dtoResponse, nil
}
