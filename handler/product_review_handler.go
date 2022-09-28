package handler

import (
	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
	"net/http"
)

func (h Handler) GetProductReviews(c *gin.Context) {
	id, err := helper.GetParam(c, "id", "int")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	reviews, err := h.productReviewService.GetProductReviews(uint64(id.(int)))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reviews)
}
