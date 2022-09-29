package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
)

func (h Handler) GetProductReviews(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	i, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	id := i.(int)

	reviews, metadata, err := h.productReviewService.GetProductReviews(uint64(id), &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"reviews":    reviews,
	})
}
