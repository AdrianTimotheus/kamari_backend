package handler

import (
	"github.com/gin-gonic/gin"
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
