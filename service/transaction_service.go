package service

import (
	"fmt"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"strconv"
	"strings"
	"time"
)

type TransactionService interface {
	CreateTransactionPayWithWallet(userModel model.User, payload dto.TransactionPayWithWalletRequest) (dto.StatusMessageResponse, error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	walletRepository      repository.WalletRepository
}

type TSConfig struct {
	TransactionRepository repository.TransactionRepository
	WalletRepository      repository.WalletRepository
}

func NewTransactionService(c TSConfig) TransactionService {
	return &transactionService{
		transactionRepository: c.TransactionRepository,
		walletRepository:      c.WalletRepository,
	}
}

func (t transactionService) CreateTransactionPayWithWallet(userModel model.User, payload dto.TransactionPayWithWalletRequest) (dto.StatusMessageResponse, error) {
	trcCode := "trc/?/wallet/?"
	trcCode = strings.Replace(trcCode, "?", time.Now().Format("01-02-2006-15-04-05"), 1)
	userIdString := strconv.Itoa(int(userModel.Id))
	trcCode = strings.Replace(trcCode, "?", "user-"+userIdString, 1)
	trcModel := model.Transaction{
		Code:            trcCode,
		PaymentMethodId: 1,
		UserId:          userModel.Id,
		IsPaid:          true,
	}

	walletModel, err := t.walletRepository.GetMyWalletDetail(userModel.Id)

	fmt.Println("this is my wallet")
	fmt.Println(walletModel.Id)

	if err != nil {
		return dto.StatusMessageResponse{Status: false}, httperror.BadRequestError("User has no active wallet", "NO_ACTIVE_WALLET")
	}

	if walletModel.Balance < payload.Amount {
		return dto.StatusMessageResponse{Status: false}, httperror.BadRequestError("Insufficient balance", "INSUFFICIENT_BALANCE")
	}

	walletTransactionModel := model.WalletTransaction{
		Amount:            payload.Amount,
		Description:       "Checkout Kamari trancation " + trcCode,
		TransactionTypeId: 2,
	}
	err = t.transactionRepository.CreateTransactionPayWithWallet(&trcModel, walletModel, &walletTransactionModel)
	if err != nil {
		return dto.StatusMessageResponse{Status: false}, httperror.InternalServerError("Internal Server Error")
	}
	redis := helper.NewRedisClient(1)
	_ = redis.Delete("stepUp-" + userModel.Email)

	return dto.StatusMessageResponse{Status: true}, nil
}
