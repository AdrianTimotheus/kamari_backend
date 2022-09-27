package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) GetCategories(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	categories, metadata, err := h.categoryService.GetCategories(query.Keyword, &query.Filters)
	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"products":   categories,
	})
}
