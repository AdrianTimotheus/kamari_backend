package handler

import (
	"fmt"
	"kamari.com/dto"
	"kamari.com/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h Handler) GetMerchantProfile(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	merchant, err := h.merchantService.GetMerchantProfile(user.Username)
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

func (h Handler) UpdateMerchantProfile(c *gin.Context) {
	var (
		req *dto.MerchantProfileReq
		err error
	)
	merchantId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	err = h.merchantService.UpdateMerchantProfile(uint64(merchantId), user.Username, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "merchant updated successfully",
	})
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
func (h Handler) GetProductsSellerAdmin(c *gin.Context) {
	var payload dto.GetMerchantProductRequest
	if err := c.ShouldBindQuery(&payload); err != nil {
		fmt.Println(err)
		_ = c.Error(httperror.BadRequestError("missing required field", "MISSING_REQUIRED_FIELD"))
		return
	}
	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}
	dtoResponse, err := h.merchantService.GetSellerAdminProducts(user.Username, payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dtoResponse)
}

func (h Handler) RegisterMerchant(c *gin.Context) {
	var (
		req *dto.MerchantRegisterRequest
		err error
	)

	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	if err = c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	err = h.merchantService.RegisterMerchant(user, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "merchant created successfully",
	})
}

func (h Handler) GetMerchantPaycheck(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	resp, err := h.merchantService.GetMerchantPaycheck(user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h Handler) WithDrawMerchantIncome(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	resp, err := h.merchantService.MerchantWithdraw(user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)

}
