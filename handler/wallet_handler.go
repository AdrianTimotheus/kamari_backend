package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
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
