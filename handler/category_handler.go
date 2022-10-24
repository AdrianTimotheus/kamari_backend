package handler

import (
	"kamari.com/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kamari.com/helper"
	"kamari.com/httperror"
)

func (h Handler) CreateCategory(c *gin.Context) {
	var payload dto.CategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	_, err := h.categoryService.AddCategory(&payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created category",
	})
}

func (h Handler) CreateSubCategory(c *gin.Context) {
	var payload dto.NewSubCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	_, err := h.categoryService.AddSubCategory(&payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created sub category",
	})
}

func (h Handler) CreateComplementaryCategory(c *gin.Context) {
	var payload dto.NewComplementaryCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	_, err := h.categoryService.AddComplementaryCategory(&payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created complementary category",
	})
}

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

func (h Handler) GetSubCategories(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	subcategories, metadata, err := h.categoryService.GetSubCategories(query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata":    metadata,
		"subcategories": subcategories,
	})
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

func (h Handler) GetComplementaryCategories(c *gin.Context) {
	var query struct {
		Keyword string `form:"keyword"`
		helper.Filters
	}

	err := c.BindQuery(&query)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_QUERY_PARAMS"))
		return
	}

	complementaries, metadata, err := h.categoryService.GetComplementaryCategories(query.Keyword, &query.Filters)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"__metadata":               metadata,
		"complementary_categories": complementaries,
	})
}

func (h Handler) GetComplementaryCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PARAM"))
		return
	}

	response, err := h.categoryService.GetComplementaryCategory(uint64(id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h Handler) UpdateCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	categoryId := id.(uint64)

	var payload dto.CategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	err := h.categoryService.UpdateCategory(categoryId, &payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated a category",
	})
}

func (h Handler) UpdateSubCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	subCategoryId := id.(uint64)

	var payload dto.EditSubCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	err := h.categoryService.UpdateSubCategory(subCategoryId, &payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated a sub category",
	})
}

func (h Handler) UpdateComplementaryCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	complementaryId := id.(uint64)

	var payload dto.EditComplementaryCategoryRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(httperror.BadRequestError(err.Error(), "INVALID_PAYLOAD"))
		return
	}
	err := h.categoryService.UpdateComplementaryCategory(complementaryId, &payload)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated a complementary category",
	})
}

func (h Handler) DeleteCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	categoryId := id.(uint64)
	err := h.categoryService.DeleteCategory(categoryId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted a category",
	})
}

func (h Handler) DeleteSubCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	subCategoryId := id.(uint64)
	err := h.categoryService.DeleteSubCategory(subCategoryId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted a sub category",
	})
}

func (h Handler) DeleteComplementaryCategory(c *gin.Context) {
	id, _ := helper.GetParam(c, "id", "uint64")
	complementaryId := id.(uint64)
	err := h.categoryService.DeleteComplementaryCategory(complementaryId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted a complementary category",
	})
}
