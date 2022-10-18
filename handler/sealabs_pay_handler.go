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

	res, err := h.sealabsAccountService.VerifySeaLabsAccount(req, &user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
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

func (h *Handler) ValidateCreateSlpAccount(c *gin.Context) {
	val := c.Query("s")

	var req *dto.SealabsPayRes

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	res, err := h.sealabsAccountService.ValidateCreateAccount(val, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	message := "successfully add an account"
	if res == nil {
		message = "transaction cancelled"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    res,
	})
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
