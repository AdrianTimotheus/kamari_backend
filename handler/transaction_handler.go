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
