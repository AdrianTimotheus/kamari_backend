package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
)

func (h Handler) GetMerchantDetail(c *gin.Context) {
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

func (h Handler) GetSellerProduct(c *gin.Context) {
	var payload dto.GetMerchantProductRequest
	if err := c.ShouldBindQuery(&payload); err != nil {
		fmt.Println(err)
		_ = c.Error(httperror.BadRequestError("missing required field", "MISSING_REQUIRED_FIELD"))
		return
	}
	res, err := helper.GetParam(c, "username", "string")
	if err != nil {
		_ = c.Error(err)
		return
	}

	username := res.(string)
	dtoResponse, err := h.merchantService.GetMerchantProducts(username, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dtoResponse)
}
