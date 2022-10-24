package handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) GetWalletDetail(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	response, err := h.walletService.GetUserWalletDetail(user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) GetWalletTransactionHistories(c *gin.Context) {
	var payload dto.WalletTransactionRequest
	err := c.ShouldBindQuery(&payload)
	if err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Query Param", "MISSING_REQUIRED_QUERY_PARAM"))
		return
	}
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	response, err := h.walletService.GetWalletTransactionHistories(user, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handler) ChangeWalletPin(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.ChangePinWalletRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Body", "MISSING_REQUIRED_BODY"))
		return
	}
	if len(payload.Password) != 6 {
		_ = c.Error(httperror.BadRequestError("Missing Required Body", "PIN should be 6 character long"))
		return
	}
	status, err := h.walletService.ChangeUserWalletPin(user.Id, payload.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.StatusMessageResponse{Status: status})
}

func (h Handler) IsUserHasWallet(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	status, err := h.walletService.IsUserHasWallet(user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.StatusMessageResponse{Status: status})
}

func (h Handler) CreateUserWallet(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.CreateWalletRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(err)
		return
	}
	if len(payload.Pin) != 6 {
		_ = c.Error(httperror.BadRequestError("Pin should 6 character long", "BAD_REQUEST"))
		return
	}
	hashedPassword := fmt.Sprintf("%X", sha256.Sum256([]byte(payload.Pin)))
	isWalletExist, err := h.walletService.IsUserHasWallet(user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if isWalletExist == true {
		_ = c.Error(httperror.BadRequestError("Wallet Already Exist", "BAD_REQUEST"))
		return
	}
	response, err := h.walletService.CreateNewWallet(user.Id, hashedPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handler) ValidatePin(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var payload dto.ValidatePinRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Body", "MISSING_REQUIRED_BODY"))
		return
	}

	stepUpToken, err := h.walletService.ValidateUserPinForPayment(user, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stepUpToken)
}

func (h Handler) TopUpWithSLP(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		return
	}

	var req *dto.PayWithSLPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	err = h.sealabsAccountService.IsValidAccountOwner(user.Id, req.CardID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, httperror.ForbiddenError())
		return
	}

	res, err := h.sealabsAccountService.TopUpWallet(req, user)

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ValidateTopUpWithSLP(c *gin.Context) {
	val := c.Query("s")

	var req *dto.SealabsPayRes

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	res, err := h.walletService.ValidateTopUpWallet(val, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	message := "top up transaction success"

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    res,
	})
}
