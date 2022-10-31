package service

import (
	"encoding/json"
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
	CreateTransactionPayWithWallet(userModel model.User, payload dto.TransactionPayWithWalletRequest) (*dto.ProcessPaymentWithWalletResponse, error)
	CreateTransactionPayWithSLP(userModel model.User, payload *dto.TransactionPaySLPRequest, signature string) (*dto.ProcessPaymentWithWalletResponse, error)
	CancelOrder(signature string) (*dto.ProcessPaymentWithWalletResponse, error)
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

func (t transactionService) CreateTransactionPayWithWallet(userModel model.User, payload dto.TransactionPayWithWalletRequest) (*dto.ProcessPaymentWithWalletResponse, error) {
	trcCode := "trc/?/wallet/?"
	trcCode = strings.Replace(trcCode, "?", time.Now().Format("01-02-2006-150405"), 1)
	userIdString := strconv.Itoa(int(userModel.Id))
	trcCode = strings.Replace(trcCode, "?", "user-"+userIdString, 1)
	trcModel := model.Transaction{
		Code:            trcCode,
		PaymentMethodId: 1,
		UserId:          userModel.Id,
		IsPaid:          true,
	}

	walletModel, err := t.walletRepository.GetMyWalletDetail(userModel.Id)

	if err != nil {
		return nil, httperror.BadRequestError("User has no active wallet", "NO_ACTIVE_WALLET")
	}

	if walletModel.Balance < payload.Amount {
		return nil, httperror.BadRequestError("Insufficient balance", "INSUFFICIENT_BALANCE")
	}

	walletTransactionModel := model.WalletTransaction{
		Amount:            payload.Amount,
		Description:       "Checkout Kamari trancation " + trcCode,
		TransactionTypeId: 2,
	}

	marketplaceVoucherModel := model.MarketplaceVoucher{ID: payload.MarketplaceVoucherID}
	var merchantVoucherModel []model.MerchantVoucher
	for _, data := range payload.MerchantVoucherIds {
		temp := model.MerchantVoucher{Id: data}
		merchantVoucherModel = append(merchantVoucherModel, temp)
	}

	var productVariantModel []model.ProductVariant
	var productPromoItems []model.ProductPromoItem
	for _, product := range payload.Products {
		temp := model.ProductPromoItem{
			StockLeft: product.Quantity,
			ProductId: product.ProductId,
		}
		productPromoItems = append(productPromoItems, temp)

		pVTempModel := model.ProductVariant{
			ID:    product.VariantId,
			Stock: int(product.Quantity),
		}
		productVariantModel = append(productVariantModel, pVTempModel)
	}

	newTrcModel, err := t.transactionRepository.CreateTransactionPayWithWallet(&trcModel, walletModel, &marketplaceVoucherModel, merchantVoucherModel, &walletTransactionModel, productPromoItems, productVariantModel)
	if err != nil {
		return nil, err
	}
	redis := helper.NewRedisClient(1)
	_ = redis.Delete("stepUp-" + userModel.Email)

	return &dto.ProcessPaymentWithWalletResponse{
		Status:        true,
		TransactionId: newTrcModel.Id,
	}, nil
}

func (t transactionService) CreateTransactionPayWithSLP(userModel model.User, payload *dto.TransactionPaySLPRequest, signature string) (*dto.ProcessPaymentWithWalletResponse, error) {
	trcCode := "trc/?/slp/?/?"
	trcCode = strings.Replace(trcCode, "?", time.Now().Format("01-02-2006-150405"), 1)
	userIdString := strconv.Itoa(int(userModel.Id))
	trcCode = strings.Replace(trcCode, "?", payload.CardNumber, 1)
	trcCode = strings.Replace(trcCode, "?", "user-"+userIdString, 1)
	trcModel := model.Transaction{
		Code:            trcCode,
		PaymentMethodId: 2,
		UserId:          userModel.Id,
		IsPaid:          true,
	}

	marketplaceVoucherModel := model.MarketplaceVoucher{ID: payload.MarketplaceVoucherID}
	var merchantVoucherModel []model.MerchantVoucher
	for _, data := range payload.MerchantVoucherIds {
		temp := model.MerchantVoucher{Id: data}
		merchantVoucherModel = append(merchantVoucherModel, temp)
	}

	var productVariantModel []model.ProductVariant
	var productPromoItems []model.ProductPromoItem
	for _, product := range payload.Products {
		temp := model.ProductPromoItem{
			StockLeft: product.Quantity,
			ProductId: product.ProductId,
		}
		productPromoItems = append(productPromoItems, temp)

		pVTempModel := model.ProductVariant{
			ID:    product.VariantId,
			Stock: int(product.Quantity),
		}
		productVariantModel = append(productVariantModel, pVTempModel)
	}

	newTrcModel, err := t.transactionRepository.CreateTransactionPayWithSLP(&trcModel, &marketplaceVoucherModel, merchantVoucherModel, productPromoItems, productVariantModel)
	if err != nil {
		return nil, err
	}

	order := dto.CreateOrderRequestSLP{
		TransactionId: newTrcModel.Id,
		Orders:        payload.Orders,
		Products:      payload.Items,
		User:          userModel,
	}

	orderJson, _ := json.Marshal(order)

	redisHelper := helper.NewRedisClient(7)
	key := redisHelper.GenerateKey("s", signature, "")

	err = redisHelper.Set(key, orderJson)

	return &dto.ProcessPaymentWithWalletResponse{
		Status:        true,
		TransactionId: newTrcModel.Id,
	}, nil
}

func (t transactionService) CancelOrder(signature string) (*dto.ProcessPaymentWithWalletResponse, error) {
	redisHelper := helper.NewRedisClient(6)
	key := redisHelper.GenerateKey("s", signature, "")
	payload, err := redisHelper.GetData(key)
	if err != nil {
		return nil, httperror.BadRequestError("Invalid Signature", "INVALID_SIGNATURE")
	}

	defer redisHelper.Delete(key)

	cancelOrder := new(dto.TransactionPaySLPRequest)
	json.Unmarshal([]byte(payload), cancelOrder)

	redisHelperSec := helper.NewRedisClient(7)
	keySec := redisHelperSec.GenerateKey("s", signature, "")
	payloadSec, err := redisHelperSec.GetData(keySec)
	if err != nil {
		return nil, httperror.BadRequestError("Invalid Signature", "INVALID_SIGNATURE")
	}

	defer redisHelperSec.Delete(keySec)

	order := new(dto.CreateOrderRequestSLP)
	json.Unmarshal([]byte(payloadSec), order)

	marketplaceVoucherModel := model.MarketplaceVoucher{ID: cancelOrder.MarketplaceVoucherID}
	var merchantVoucherModel []model.MerchantVoucher
	for _, data := range cancelOrder.MerchantVoucherIds {
		temp := model.MerchantVoucher{Id: data}
		merchantVoucherModel = append(merchantVoucherModel, temp)
	}

	var productVariantModel []model.ProductVariant
	var productPromoItems []model.ProductPromoItem
	for _, product := range cancelOrder.Products {
		temp := model.ProductPromoItem{
			StockLeft: product.Quantity,
			ProductId: product.ProductId,
		}
		productPromoItems = append(productPromoItems, temp)

		pVTempModel := model.ProductVariant{
			ID:    product.VariantId,
			Stock: int(product.Quantity),
		}
		productVariantModel = append(productVariantModel, pVTempModel)
	}

	_, err = t.transactionRepository.CancelCreateTransactionPayWithSLP(order, &marketplaceVoucherModel, merchantVoucherModel, productPromoItems, productVariantModel)
	if err != nil {
		return nil, err
	}

	return nil, err
}
