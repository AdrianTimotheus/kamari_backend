package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
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

	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "BAD_PARAMS"))
		return
	}

	accountId := p.(uint64)
	if err := h.sealabsAccountService.SetDefaultSealabsAccount(user.Id, accountId); err != nil {
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
		_ = c.Error(httperror.UnauthorizedError())
		return
	}

	accounts, err := h.sealabsAccountService.GetUsesSealabsAccounts(user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h Handler) CheckSealabsOwner(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "BAD_PARAMS"))
		return
	}

	accountId := p.(uint64)
	err = h.sealabsAccountService.IsValidAccountOwner(user.Id, accountId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, httperror.ForbiddenError())
		return
	}
}

func (h Handler) DeleteSealabsAccount(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "BAD_PARAMS"))
		return
	}

	accountId := p.(uint64)
	err = h.sealabsAccountService.RemoveAccount(accountId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully deleted an account",
	})
}
