package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
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
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata": metadata,
		"categories": categories,
	})
}

func (h Handler) GetCategory(c *gin.Context) {
	categoryId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	category, err := h.categoryService.GetCategory(uint64(categoryId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h Handler) GetSubcategories(c *gin.Context) {
	response, err := h.categoryService.GetSubcategories()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) GetSubcategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	response, err := h.categoryService.GetSubcategory(uint64(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) GetComplementaries(c *gin.Context) {
	response, err := h.categoryService.GetComplementaries()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) GetComplementary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	response, err := h.categoryService.GetComplementary(uint64(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
