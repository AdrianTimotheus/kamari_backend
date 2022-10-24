package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/dto"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) GetAllCourier(c *gin.Context) {
	resp, err := h.courierService.GetAllCouriers()
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetShippingCost(c *gin.Context) {
	var payloads dto.GetShippingCost
	if err := c.ShouldBindJSON(&payloads); err != nil {
		_ = c.Error(httperror.BadRequestError("missing required field", "MISSING_REQUIRED_FIELD"))
		return
	}
	res := h.courierService.GetShippingCost(&payloads)

	c.JSON(http.StatusOK, res)
}
