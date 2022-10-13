package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"time"
)

type WalletService interface {
	GetUserWalletDetail(userModel model.User) (*dto.DetailWalletResponse, error)
	GetWalletTransactionHistories(userModel model.User, payload dto.WalletTransactionRequest) (*dto.WalletTransactionResponse, error)
	ChangeUserWalletPin(userId uint64, Password string) (bool, error)
	IsUserHasWallet(userId uint64) (bool, error)
	CreateNewWallet(userId uint64, hashedPassword string) (*dto.DetailWalletResponse, error)
	ValidateUserPinForPayment(user model.User, payload dto.ValidatePinRequest) (*dto.TokenResponse, error)
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

func (w walletService) IsUserHasWallet(userId uint64) (bool, error) {
	count, err := w.walletRepository.CountUserWallet(userId)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (w walletService) CreateNewWallet(userId uint64, hashedPassword string) (*dto.DetailWalletResponse, error) {
	wModel := model.Wallet{
		Id:      707000 + userId,
		Balance: 0,
		Pin:     hashedPassword,
		UserId:  userId,
	}
	nWModel, err := w.walletRepository.CreateWallet(wModel)
	if err != nil {
		return nil, err
	}
	dtoResponse := dto.DetailWalletResponse{
		Id:      nWModel.Id,
		Balance: nWModel.Balance,
	}
	return &dtoResponse, nil
}

type levelOneJWTToken struct {
	jwt.RegisteredClaims
	User  *model.User `json:"user"`
	Scope string      `json:"scope"`
}

func generateStepUpToken(user *model.User) (*dto.TokenResponse, error) {
	idExp := config.Config.JWTExpireInMinutes * 60
	unixTime := time.Now().Unix()
	tokenExp := unixTime + int64(idExp)
	timeExpire := jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	timeNow := jwt.NumericDate{Time: time.Now()}
	claims := &levelOneJWTToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.AppName,
			IssuedAt:  &timeNow,
			ExpiresAt: &timeExpire,
		},
		User:  user,
		Scope: "Level1",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Config.JWTSecret)

	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{IDToken: tokenString}, err
}

func (w walletService) ValidateUserPinForPayment(user model.User, payload dto.ValidatePinRequest) (*dto.TokenResponse, error) {
	hashedPin := fmt.Sprintf("%X", sha256.Sum256([]byte(payload.Pin)))
	fmt.Println(user.Id)
	fmt.Println(hashedPin)
	if !w.walletRepository.ValidatePin(user.Id, hashedPin) {
		return nil, httperror.ForbiddenError()
	}

	dtoResponse, err := generateStepUpToken(&user)
	if err != nil {
		return nil, httperror.InternalServerError("Internal Server Error")
	}

	redis := helper.NewRedisClient(1)
	stepUpKey := "stepUp-" + user.Email

	err = redis.SetData(stepUpKey, dtoResponse.IDToken)
	if err != nil {
		return nil, httperror.InternalServerError("Internal Server Error")
	}

	return dtoResponse, nil
}
