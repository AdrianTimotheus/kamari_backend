package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetCities(c *gin.Context) {
	dataResp, err := h.cityService.GetListCityPerProvince()
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(http.StatusOK, dataResp)
}
