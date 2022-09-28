package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"time"
)

type UserService interface {
	CreateNewUser(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	CheckAuthPassword(payload dto.CheckAuthPassword, userId uint64) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

type USConfig struct {
	UserRepository repository.UserRepository
}

func NewUserService(c USConfig) UserService {
	return userService{
		userRepository: c.UserRepository,
	}
}

type idTokenClaims struct {
	jwt.RegisteredClaims
	User *model.User `json:"user"`
}

func (u *userService) generateJWTToken(user *model.User) (*dto.TokenResponse, error) {
	idExp := config.Config.JWTExpireInMinutes * 5
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
