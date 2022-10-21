package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) ValidateCurrentSellerVoucher(c *gin.Context) {
	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	v, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, httperror.BadRequestError(err.Error(), "INVALID_ID_PARAMS"))
		return
	}

	vouchedId := v.(uint64)
	err = h.merchantVoucherService.ValidateVoucherOwner(user.Username, vouchedId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, httperror.ForbiddenError())
	}
}

func (h Handler) DeleteMerchantVoucher(c *gin.Context) {
	v, _ := helper.GetParam(c, "id", "uint64")
	voucherId := v.(uint64)

	err := h.merchantVoucherService.DeleteVoucher(voucherId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted a voucher",
	})
}

func (h Handler) UpdateMerchantVoucher(c *gin.Context) {
	v, _ := helper.GetParam(c, "id", "uint64")
	voucherId := v.(uint64)

	var payload dto.MerchantVoucherUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	err := h.merchantVoucherService.UpdateVoucher(voucherId, &payload)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated a voucher",
	})
}

func (h Handler) GetMerchantVouchers(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	u, err := h.merchantService.GetMerchantDetail(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	vouchers, metadata, err := h.merchantVoucherService.GetVouchers(u.Id, query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"vouchers":   vouchers,
	})
}

func (h Handler) CreateMerchantVoucher(c *gin.Context) {
	var payload dto.MerchantVoucherRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}

	user, err := helper.GetUserContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	err = h.merchantVoucherService.CreateVoucher(&payload, user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created voucher",
	})
}

func (h Handler) GetMerchantVoucherDetail(c *gin.Context) {
	p, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = httperror.BadRequestError(err.Error(), "BAD_PARAMS")
		return
	}

	voucherId := p.(uint64)
	voucher, err := h.merchantVoucherService.GetVoucherDetail(voucherId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, voucher)
}

func (h Handler) GetMerchantVoucherByMerchantId(c *gin.Context) {
	v, _ := helper.GetParam(c, "id", "uint64")
	merchantId := v.(uint64)
	resp, err := h.merchantVoucherService.GetAllActiveVoucher(merchantId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
