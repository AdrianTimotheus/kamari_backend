package handler

import (
	"errors"
	"kamari.com/service"
)

var (
	ErrFailedBindingValidation = errors.New("failed binding validation")
)

type Handler struct {
	userService     service.UserService
	provinceService service.ProvinceService
}

type HConfig struct {
	UserService     service.UserService
	ProvinceService service.ProvinceService
}

func New(c HConfig) Handler {
	return Handler{
		userService:     c.UserService,
		provinceService: c.ProvinceService,
	}
}
