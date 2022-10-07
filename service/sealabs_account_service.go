package service

import (
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type SealabsAccountService interface {
	CreateNewAccount(*dto.SealabsAccountRequest, *model.User) error
	GetUsesSealabsAccounts(userId uint64) ([]dto.SealabsPayAccountResponse, error)
	SetDefaultSealabsAccount(userId, accountId uint64) error
}

type sealabsAccountService struct {
	sealabsAccountRepository repository.SealabsAccountRepository
}

type SASConfig struct {
	SealabsAccountRepository repository.SealabsAccountRepository
}

func NewSealabsAccountService(c SASConfig) SealabsAccountService {
	return sealabsAccountService{
		sealabsAccountRepository: c.SealabsAccountRepository,
	}
}

func (s sealabsAccountService) CreateNewAccount(req *dto.SealabsAccountRequest, user *model.User) error {
	sealabsAccount := &model.SealabsPayAccount{
		Title:      req.Title,
		CardNumber: req.CardNumber,
		UserId:     user.Id,
	}

	if s.sealabsAccountRepository.HasSameCard(sealabsAccount) {
		return httperror.BadRequestError("Card with given data title already exist", "DUPLICATE_CARD_FOUND")
	}

	accounts, err := s.sealabsAccountRepository.GetAccounts(sealabsAccount.UserId)
	if err != nil {
		return err
	}

	accountId, err := s.sealabsAccountRepository.CreateAccount(sealabsAccount)
	if err != nil {
		return err
	}

	if req.IsDefault || len(accounts) < 1 {
		err := s.sealabsAccountRepository.ChangeDefaultTransaction(sealabsAccount.UserId, accountId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s sealabsAccountService) GetUsesSealabsAccounts(userId uint64) ([]dto.SealabsPayAccountResponse, error) {
	accounts, err := s.sealabsAccountRepository.GetAccounts(userId)
	if err != nil {
		return nil, err
	}

	response := make([]dto.SealabsPayAccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = dto.SealabsPayAccountResponse{
			Id:         account.Id,
			Title:      account.Title,
			CardNumber: account.CardNumber,
			IsDefault:  account.IsDefault,
		}
	}

	return response, nil
}

func (s sealabsAccountService) SetDefaultSealabsAccount(userId, accountId uint64) error {
	err := s.sealabsAccountRepository.ChangeDefaultTransaction(userId, accountId)
	if err != nil {
		return err
	}

	return nil
}
