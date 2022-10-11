package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgconn"
	"github.com/mailjet/mailjet-apiv3-go"
	"gorm.io/gorm"
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/middleware"
	"kamari.com/model"
	"kamari.com/repository"
)

type UserService interface {
	CreateNewUser(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	LoginWithGmail(req *dto.GoogleAuthTokenReq) (*dto.TokenResponse, error)
	CheckAuthPassword(payload dto.CheckAuthPassword, userId uint64) (bool, error)
	//GenerateOTP(user model.User) (dto.StatusMessageResponse, error)
	GenerateOTPRedis(user model.User) (dto.StatusMessageResponse, error)
	AuthOTP(user model.User, OTPCode dto.AuthOTPRequest) (dto.StatusMessageResponse, error)
	GetUserDetail(userId uint64) (*model.User, error)
	EditUserDetail(req dto.EditUserDetailReq, userId uint64) (*model.User, error)
}

type userService struct {
	userRepository    repository.UserRepository
	otpCodeRepository repository.OtpCodeRepository
}

type USConfig struct {
	UserRepository    repository.UserRepository
	OtpCodeRepository repository.OtpCodeRepository
}

func NewUserService(c USConfig) UserService {
	return userService{
		userRepository:    c.UserRepository,
		otpCodeRepository: c.OtpCodeRepository,
	}
}

type idTokenClaims struct {
	jwt.RegisteredClaims
	User *model.User `json:"user"`
}

func (u userService) generateJWTToken(user *model.User) (*dto.TokenResponse, error) {
	idExp := config.Config.JWTExpireInMinutes * 60
	unixTime := time.Now().Unix()
	tokenExp := unixTime + int64(idExp)
	timeExpire := jwt.NumericDate{Time: time.Unix(tokenExp, 0)}
	timeNow := jwt.NumericDate{Time: time.Now()}
	claims := &idTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Config.AppName,
			IssuedAt:  &timeNow,
			ExpiresAt: &timeExpire,
		},
		User: user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Config.JWTSecret)

	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{IDToken: tokenString}, err
}

func (u userService) CreateNewUser(req *dto.RegisterRequest) error {
	if helper.EmailNotValid(req.Username, req.Password) {
		return httperror.BadRequestError("password cannot contains username", "INVALID_PASSWORD_PAYLOAD")
	}

	user := req.ToModel(req)

	hash := fmt.Sprintf("%X", sha256.Sum256([]byte(req.Password)))
	user.Password = hash

	err := u.userRepository.CreateUser(user)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr):
			return httperror.InternalServerError("duplicate key values for unique constraint")
		default:
			return httperror.InternalServerError(err.Error())
		}
	}

	return nil
}

func (u userService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := u.userRepository.MatchingCredential(req.Email, req.Password)
	if err != nil || user == nil {
		return nil, httperror.InvalidLogin("Invalid email or password")
	}
	token, err := u.generateJWTToken(user)
	return token, err
}

func (u userService) LoginWithGmail(req *dto.GoogleAuthTokenReq) (*dto.TokenResponse, error) {
	unauthorizedErr := httperror.UnauthorizedError()
	inputToken, err := middleware.VerifyGoogleToken(req.Token)
	if err != nil {
		return nil, unauthorizedErr
	}

	user, err := u.userRepository.FindUserByEmail(inputToken.Email)
	if err != nil || user == nil {
		return nil, httperror.InvalidLogin("Invalid email")
	}
	token, err := u.generateJWTToken(user)
	return token, err
}

func (u userService) CheckAuthPassword(payload dto.CheckAuthPassword, userId uint64) (bool, error) {
	hashedPassword := fmt.Sprintf("%X", sha256.Sum256([]byte(payload.Password)))
	_, err := u.userRepository.CheckPassword(userId, hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return false, httperror.BadRequestError(err.Error(), "NO_ROWS")
		default:
			return false, httperror.InternalServerError(err.Error())
		}
	}
	return true, nil
}

//func (u userService) GenerateOTP(user model.User) (dto.StatusMessageResponse, error) {
//	fmt.Println("here", user.Email)
//	if _, err := u.otpCodeRepository.ClearUserOTP(user.Id); err != nil {
//		return dto.StatusMessageResponse{Status: false}, err
//	}
//
//	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
//	rand.Seed(time.Now().UnixNano())
//	otpCode := make([]rune, 6)
//	for i := range otpCode {
//		otpCode[i] = letterRunes[rand.Intn(len(letterRunes))]
//	}
//	otpString := string(otpCode)
//
//	otpModel := model.OtpCode{
//		Code:   otpString,
//		UserId: user.Id,
//	}
//
//	_, err := u.otpCodeRepository.CreateOTP(otpModel)
//	if err != nil {
//		return dto.StatusMessageResponse{Status: false}, err
//	}
//	err = SendOTPMailToUser(user, otpString)
//	if err != nil {
//		return dto.StatusMessageResponse{Status: false}, err
//	}
//	return dto.StatusMessageResponse{Status: true}, nil
//}

func (u userService) GenerateOTPRedis(user model.User) (dto.StatusMessageResponse, error) {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rand.Seed(time.Now().UnixNano())
	otpCode := make([]rune, 6)
	for i := range otpCode {
		otpCode[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	otpString := string(otpCode)
	otpKey := "OTP" + strconv.Itoa(int(user.Id))
	redisHelper := helper.NewRedisClient()
	hashedOtpString := fmt.Sprintf("%X", sha256.Sum256([]byte(otpString)))
	err := redisHelper.SetData(otpKey, hashedOtpString)
	if err != nil {
		return dto.StatusMessageResponse{Status: false}, err
	}
	err = SendOTPMailToUser(user, otpString)
	if err != nil {
		return dto.StatusMessageResponse{Status: false}, err
	}
	return dto.StatusMessageResponse{Status: true}, nil
}

func SendOTPMailToUser(userModel model.User, otpCode string) error {
	mailjetClient := mailjet.NewMailjetClient("b0cc175b3bbae4dcc34d24d4cb3268df", "e828d182d7aa0a62b2fbc728e7ee440f")
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "kamari.commerce@gmail.com",
				Name:  "Kamari Commerce [No Reply]",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: userModel.Email,
					Name:  userModel.Name,
				},
			},
			Subject:  "OTP for required action [DO NOT SHARE TO ANYONE].",
			TextPart: "Email OTP Kamari Commerce",
			HTMLPart: "<h3>Dear Customer , Here is your OTP for further required action</h3><br/>" +
				"<h1>" + otpCode + "</h1>" +
				"<br/><h2>PLEASE DON'T SHARE THIS CODE TO ANYONE</h2>",
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}

func (u userService) AuthOTP(user model.User, OTPCode dto.AuthOTPRequest) (dto.StatusMessageResponse, error) {
	otpKey := "OTP" + strconv.Itoa(int(user.Id))
	redisHelper := helper.NewRedisClient()
	hashedOtpCode := fmt.Sprintf("%X", sha256.Sum256([]byte(OTPCode.Code)))
	res, err := redisHelper.GetData(otpKey)
	if err != nil {
		return dto.StatusMessageResponse{Status: false}, nil
	}
	if res != hashedOtpCode {
		return dto.StatusMessageResponse{Status: false}, nil
	}
	return dto.StatusMessageResponse{Status: true}, nil
}

func (u userService) GetUserDetail(userId uint64) (*model.User, error) {
	response, err := u.userRepository.GetUserDetail(userId)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (u userService) EditUserDetail(req dto.EditUserDetailReq, userId uint64) (*model.User, error) {
	user := model.User{
		Id:          userId,
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
		PicUrl:      req.PicUrl,
		GenderId:    req.GenderId,
	}

	response, err := u.userRepository.EditUserDetail(user)

	if err != nil {
		return nil, err
	}
	return response, err
}
