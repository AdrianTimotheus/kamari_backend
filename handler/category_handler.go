package handler

import (
	"net/http"
	"strconv"

	"kamari.com/dto"
	"kamari.com/model"

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

func (h Handler) GetSellerCategories(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	params := c.Request.URL.Query()

	response, totalPage, err := h.categoryService.GetSellerCategories(user.Id, params)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"seller_categories": response,
		"total_page":        totalPage,
	})
}

func (h Handler) GetSellerPageCategories(c *gin.Context) {
	res, err := helper.GetParam(c, "username", "string")
	if err != nil {
		_ = c.Error(err)
		return
	}
	username := res.(string)

	response, err := h.categoryService.GetSellerPageCategories(username)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"seller_categories": response,
	})
}

func (h Handler) UpdateSellerCategory(c *gin.Context) {
	var req dto.SellerCategoryRequest

	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.categoryService.UpdateSellerCategory(id.(uint64), user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seller category updated",
	})
}

func (h Handler) CreateSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var req dto.SellerCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	err := h.categoryService.CreateSellerCategory(user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seller category created",
	})
}

func (h Handler) DeleteSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.categoryService.DeleteSellerCategory(user.Id, id.(uint64))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seller category deleted",
	})
}

func (h Handler) GetSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	title, response, err := h.categoryService.GetSellerCategory(user.Id, id.(uint64))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"title": title,
		"items": response,
	})
}

func (h Handler) GetAvailableProducts(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)
	params := c.Request.URL.Query()

	response, totalPage, err := h.categoryService.GetAvailableProducts(user.Id, params)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      response,
		"total_page": totalPage,
	})
}

func (h Handler) AddProductsToSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var req dto.AddRemoveProductFromCategory

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	id, err := helper.GetParam(c, "id", "uint64")
	if err != nil {
		_ = c.Error(httperror.BadRequestError(helper.ErrInvalidReturnType.Error(), "BAD_PARAMS"))
		return
	}

	err = h.categoryService.AddProductToSellerCategory(user.Id, req, id.(uint64))

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to seller category",
	})
}

func (h Handler) RemoveProductFromSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	var req dto.AddRemoveProductFromCategory

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(httperror.BadRequestError(ErrFailedBindingValidation.Error(), "INVALID_PAYLOAD"))
		return
	}

	err := h.categoryService.DeleteProductFromSellerCategory(user.Id, req)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from seller category",
	})
}

func (h Handler) RemoveSpecificProductFromSellerCategory(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(model.User)

	id, _ := helper.GetParam(c, "id", "uint64")
	productId := id.(uint64)

	err := h.categoryService.RemoveSpecificProductFromSellerCategory(user.Id, productId)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from seller category",
	})
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
