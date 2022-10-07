package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) AddSealabsPayAccount(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req *dto.SealabsAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	err = h.sealabsAccountService.CreateNewAccount(req, &user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success created new sealabs account",
	})
}

func (h Handler) SetDefaultSealabsAccount(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.SealabsAccountChangeDefault
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(ErrFailedBindingValidation)
		return
	}

	if err := h.sealabsAccountService.SetDefaultSealabsAccount(user.Id, req.SealabsAccountId); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully changed sealabspay account default",
	})
}

func (h Handler) GetUserSealabsAccount(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	accounts, err := h.sealabsAccountService.GetUsesSealabsAccounts(user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}
