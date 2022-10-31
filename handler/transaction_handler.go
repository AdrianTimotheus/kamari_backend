package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) PayCheckoutWithWallet(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var payload dto.TransactionPayWithWalletRequest
	if err := c.ShouldBind(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Field", "MISSING_REQUIRED_FIELD"))
		return
	}
	resp, err := h.transactionService.CreateTransactionPayWithWallet(user, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h Handler) PayCheckoutWithSLP(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	var req *dto.TransactionPaySLPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	err := h.sealabsAccountService.IsValidAccountOwner(user.Id, req.CardID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, httperror.ForbiddenError())
		return
	}
	req.User = user

	res, signature, err := h.sealabsAccountService.PayOrder(req, user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp, err := h.transactionService.CreateTransactionPayWithSLP(user, req, signature)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"transaction_data": resp,
		"redirect_url":     res,
	})
}

func (h Handler) GetWithdrawHistories(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	resp, err := h.walletTransactionService.GetMerchantWithdrawHistories(user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
