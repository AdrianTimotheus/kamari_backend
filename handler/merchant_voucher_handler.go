package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) ValidateCurrentSeller(c *gin.Context) {
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
