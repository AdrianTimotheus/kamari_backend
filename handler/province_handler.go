package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetAllProvinces(c *gin.Context) {
	response, err := h.provinceService.GetAllProvinces()

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
