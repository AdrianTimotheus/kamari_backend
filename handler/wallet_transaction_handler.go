package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"net/http"
)

func (h Handler) GetWalletTransactionDetail(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var payload dto.GetDetailWalletTransactionRequest
	if err := c.ShouldBind(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError("Missing Required Query Param", "MISSING_REQUIRED_QUERY_PARAM"))
		return
	}

	response, err := h.walletTransactionService.GetDetailTransaction(payload.Id, user.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}
