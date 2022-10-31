package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type SealabsAccountService interface {
	CreateNewAccount(req *dto.SealabsAccountRequest, user *model.User) (*model.SealabsPayAccount, error)
	GetUsesSealabsAccounts(userId uint64) ([]dto.SealabsPayAccountResponse, error)
	SetDefaultSealabsAccount(userId, accountId uint64) error
	IsValidAccountOwner(userId, accountId uint64) error
	RemoveAccount(accountId uint64) error
	VerifySeaLabsAccount(req *dto.SealabsAccountRequest, user *model.User) (*dto.SLPPayRes, error)
	ValidateCreateAccount(sign string, req *dto.SealabsPayRes) (*dto.SealabsPayAccountResponse, error)
	PostSeaLabsPay(req string) (*string, error)
	TopUpWallet(payload *dto.PayWithSLPReq, user model.User) (*dto.SLPPayRes, error)
	generateCallBackUrl(sign string, endPoint string) string
	PayOrder(payload *dto.TransactionPaySLPRequest, user model.User) (*dto.SLPPayRes, string, error)
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

func (s sealabsAccountService) CreateNewAccount(req *dto.SealabsAccountRequest, user *model.User) (*model.SealabsPayAccount, error) {
	accounts, err := s.sealabsAccountRepository.GetAccounts(user.Id)
	if err != nil {
		return nil, err
	}

	val := req.IsDefault
	if len(accounts) < 1 {
		val = true
	}

	sealabsAccount := &model.SealabsPayAccount{
		Title:      req.Title,
		CardNumber: req.CardNumber,
		UserId:     user.Id,
		IsDefault:  val,
	}

	if s.sealabsAccountRepository.HasSameCard(sealabsAccount) {
		return nil, httperror.BadRequestError("Card with given data title already exist", "DUPLICATE_CARD_FOUND")
	}

	return sealabsAccount, nil
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

func (s sealabsAccountService) IsValidAccountOwner(userId, accountId uint64) error {
	valid, err := s.sealabsAccountRepository.IsValidOwner(userId, accountId)
	if err != nil || !valid {
		return httperror.ForbiddenError()
	}

	return nil
}

func (s sealabsAccountService) RemoveAccount(accountId uint64) error {
	account := model.SealabsPayAccount{Id: accountId}
	err := s.sealabsAccountRepository.RemoveAccountTransaction(&account)
	if err != nil {
		return err
	}

	return nil
}

func (s sealabsAccountService) VerifySeaLabsAccount(req *dto.SealabsAccountRequest, user *model.User) (*dto.SLPPayRes, error) {
	slpAcc, err := s.CreateNewAccount(req, user)
	if err != nil {
		return nil, err
	}
	slpJson, _ := json.Marshal(slpAcc)

	const signatureLen = 10
	callBackSignature := helper.RandomString(signatureLen)
	callBackUrl := s.generateCallBackUrl(callBackSignature, "users/sealabspay/validate")
	redirectUrl := "http://localhost:3000/payment"

	postReq := dto.SLPPayReq{
		CardNumber:   req.CardNumber,
		Amount:       0,
		MerchantCode: config.Config.SeaLabsPayConfig.MerchantCode,
		RedirectUrl:  redirectUrl,
		CallbackUrl:  callBackUrl,
	}
	helper.CreateSignature(config.Config.SeaLabsPayConfig.ApiKey, &postReq)

	res, err := s.PostSeaLabsPay(helper.Encode(&postReq))

	if err != nil {
		return nil, err
	}

	r := new(dto.SLPPayRes)

	redisHelper := helper.NewRedisClient(3)
	key := redisHelper.GenerateKey("s", callBackSignature, "")

	err = redisHelper.Set(key, slpJson)

	r = r.Create(*res)

	return r, err
}

func (s sealabsAccountService) ValidateCreateAccount(sign string, req *dto.SealabsPayRes) (*dto.SealabsPayAccountResponse, error) {
	redisHelper := helper.NewRedisClient(3)
	key := redisHelper.GenerateKey("s", sign, "")
	payload, err := redisHelper.GetData(key)
	if err != nil {
		return nil, httperror.BadRequestError("Invalid Signature", "INVALID_SIGNATURE")
	}

	defer redisHelper.Delete(key)

	if req.Status == "TXN_FAILED" {
		return nil, err
	}

	slp := new(model.SealabsPayAccount)
	json.Unmarshal([]byte(payload), slp)

	accountId, err := s.sealabsAccountRepository.CreateAccount(slp)
	if err != nil {
		return nil, err
	}

	if slp.IsDefault {
		err := s.sealabsAccountRepository.ChangeDefaultTransaction(slp.UserId, accountId)
		if err != nil {
			return nil, err
		}
	}

	res := new(dto.SealabsPayAccountResponse)
	if err == nil {
		res = res.CreateFromSeaLabsPay(slp)
	}

	return res, err
}

func (s sealabsAccountService) TopUpWallet(payload *dto.PayWithSLPReq, user model.User) (*dto.SLPPayRes, error) {
	topUpReq := dto.PayWithSLPReq{
		UserID:     user.Id,
		CardNumber: payload.CardNumber,
		Amount:     payload.Amount,
	}

	reqJson, _ := json.Marshal(topUpReq)

	const signatureLen = 10
	callBackSignature := helper.RandomString(signatureLen)
	callBackUrl := s.generateCallBackUrl(callBackSignature, "wallets/topup")
	redirectUrl := "http://localhost:3000/wallet"

	postReq := dto.SLPPayReq{
		CardNumber:   payload.CardNumber,
		Amount:       payload.Amount,
		MerchantCode: config.Config.SeaLabsPayConfig.MerchantCode,
		RedirectUrl:  redirectUrl,
		CallbackUrl:  callBackUrl,
	}
	helper.CreateSignature(config.Config.SeaLabsPayConfig.ApiKey, &postReq)

	res, err := s.PostSeaLabsPay(helper.Encode(&postReq))

	if err != nil {
		return nil, err
	}

	r := new(dto.SLPPayRes)

	redisHelper := helper.NewRedisClient(4)
	key := redisHelper.GenerateKey("s", callBackSignature, "")

	err = redisHelper.Set(key, reqJson)

	r = r.Create(*res)

	return r, err
}

func (s sealabsAccountService) PayOrder(payload *dto.TransactionPaySLPRequest, user model.User) (*dto.SLPPayRes, string, error) {
	reqJson, _ := json.Marshal(payload)

	const signatureLen = 10
	callBackSignature := helper.RandomString(signatureLen)
	callBackUrl := s.generateCallBackUrl(callBackSignature, "checkout/proceed/slp")
	redirectUrl := "http://localhost:3000/payment"

	postReq := dto.SLPPayReq{
		CardNumber:   payload.CardNumber,
		Amount:       int(payload.Amount),
		MerchantCode: config.Config.SeaLabsPayConfig.MerchantCode,
		RedirectUrl:  redirectUrl,
		CallbackUrl:  callBackUrl,
	}
	helper.CreateSignature(config.Config.SeaLabsPayConfig.ApiKey, &postReq)

	res, err := s.PostSeaLabsPay(helper.Encode(&postReq))

	if err != nil {
		return nil, "", err
	}

	r := new(dto.SLPPayRes)

	redisHelper := helper.NewRedisClient(6)
	key := redisHelper.GenerateKey("s", callBackSignature, "")

	err = redisHelper.Set(key, reqJson)

	r = r.Create(*res)

	return r, callBackSignature, err
}

func (s sealabsAccountService) generateCallBackUrl(sign string, endPoint string) string {
	//Change when deployed
	baseUrl := "http://api.kamarioddspace.com/api/" + endPoint
	callBackUrl := fmt.Sprintf("%s?s=%s", baseUrl, sign)
	return callBackUrl
}

func (s sealabsAccountService) PostSeaLabsPay(req string) (*string, error) {
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	resp, err := client.Post("https://slp.air-sipp.com/api/v1/transaction/pay", "application/x-www-form-urlencoded", strings.NewReader(req))
	if err != nil {
		return nil, httperror.InternalServerError("Can't connect to SeaLabs Pay")
	}

	l, err := resp.Location()
	redirectUrl := fmt.Sprintf("%v", l)
	if err != nil {
		return nil, httperror.BadRequestError("Invalid Account", "INVALID_ACCOUNT")
	}

	return &redirectUrl, err
}
