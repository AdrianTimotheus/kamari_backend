package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
)

func (h Handler) GetSellerDetail(c *gin.Context) {
	res, err := helper.GetParam(c, "username", "string")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	username := res.(string)
	merchant, err := h.merchantService.GetMerchantDetail(username)
	if err != nil {
		e := err.(httperror.AppError)
		_ = c.Error(httperror.AppError{
			StatusCode: e.StatusCode,
			Code:       e.Code,
			Message:    e.Message,
		})
		return
	}

	c.JSON(http.StatusOK, merchant)
}
