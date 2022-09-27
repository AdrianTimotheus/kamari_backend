package service

import (
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/repository"
)

type CategoryService interface {
	GetCategories(keyword string, filters *helper.Filters) ([]dto.CategoryResponse, *helper.Metadata, error)
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

type CaSConfig struct {
	CategoryRepository repository.CategoryRepository
}

func NewCategoryService(c CaSConfig) CategoryService {
	return categoryService{
		categoryRepository: c.CategoryRepository,
	}
}

func (c categoryService) GetCategories(keyword string, filters *helper.Filters) ([]dto.CategoryResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)

	categories, metadata, err := c.categoryRepository.GetCategories(keyword, filters)
	if err != nil {
		return nil, nil, httperror.InternalServerError(err.Error())
	}

	categoryResponse := make([]dto.CategoryResponse, len(categories))
	for i, category := range categories {
		c := dto.CategoryResponse{
			Id:            category.Id,
			Name:          category.Name,
			PicUrl:        category.PicUrl,
			Subcategories: category.Subcategories,
		}
		categoryResponse[i] = c
	}

	return categoryResponse, metadata, nil
}
