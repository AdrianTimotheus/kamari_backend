package service

import (
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
)

type CategoryService interface {
	AddCategory(req *dto.CategoryRequest) (*model.Category, error)
	AddSubCategory(req *dto.NewSubCategoryRequest) (*model.Subcategory, error)
	AddComplementaryCategory(req *dto.NewComplementaryCategoryRequest) (*model.Complementary, error)
	GetCategory(categoryId uint64) (*model.Category, error)
	GetCategories(keyword string, filters *helper.Filters) ([]dto.CategoryResponse, *helper.Metadata, error)
	GetSubCategories(keyword string, filters *helper.Filters) ([]dto.SubCategoryResponse, *helper.Metadata, error)
	GetSubcategory(id uint64) (*model.Subcategory, error)
	GetComplementaryCategories(keyword string, filters *helper.Filters) ([]dto.ComplementaryCategoryResponse, *helper.Metadata, error)
	GetComplementaryCategory(id uint64) (*model.Complementary, error)
	UpdateCategory(categoryId uint64, req *dto.CategoryRequest) error
	UpdateSubCategory(subCategoryId uint64, req *dto.EditSubCategoryRequest) error
	UpdateComplementaryCategory(complementaryId uint64, req *dto.EditComplementaryCategoryRequest) error
	DeleteCategory(categoryId uint64) error
	DeleteSubCategory(subCategoryId uint64) error
	DeleteComplementaryCategory(complementaryId uint64) error
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

func (c categoryService) AddCategory(req *dto.CategoryRequest) (*model.Category, error) {
	return c.categoryRepository.SaveCategory(&model.Category{
		Name:   req.Name,
		PicUrl: req.PicUrl,
	})
}

func (c categoryService) AddSubCategory(req *dto.NewSubCategoryRequest) (*model.Subcategory, error) {
	return c.categoryRepository.SaveSubCategory(&model.Subcategory{
		Name:       req.Name,
		CategoryId: req.CategoryId,
	})
}

func (c categoryService) AddComplementaryCategory(req *dto.NewComplementaryCategoryRequest) (*model.Complementary, error) {
	return c.categoryRepository.SaveComplementaryCategory(&model.Complementary{
		Name:          req.Name,
		SubcategoryId: req.SubCategoryId,
	})
}

func (c categoryService) GetCategory(categoryId uint64) (*model.Category, error) {
	return c.categoryRepository.GetCategory(&model.Category{Id: categoryId})
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

func (c categoryService) GetSubCategories(keyword string, filters *helper.Filters) ([]dto.SubCategoryResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)

	subcategories, metadata, err := c.categoryRepository.GetSubCategories(keyword, filters)
	if err != nil {
		return nil, nil, httperror.InternalServerError(err.Error())
	}

	subcategoryResponse := make([]dto.SubCategoryResponse, len(subcategories))
	for i, subcategory := range subcategories {
		c := dto.SubCategoryResponse{
			Id:   subcategory.Id,
			Name: subcategory.Name,
		}
		subcategoryResponse[i] = c
	}
	return subcategoryResponse, metadata, nil
}

func (c categoryService) GetSubcategory(id uint64) (*model.Subcategory, error) {
	return c.categoryRepository.GetSubcategory(id)
}

func (c categoryService) GetComplementaryCategories(keyword string, filters *helper.Filters) ([]dto.ComplementaryCategoryResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 30)
	filters.Page = helper.ReadInt(filters.Page, 1)

	complementaries, metadata, err := c.categoryRepository.GetComplementaryCategories(keyword, filters)
	if err != nil {
		return nil, nil, httperror.InternalServerError(err.Error())
	}

	complementaryResponse := make([]dto.ComplementaryCategoryResponse, len(complementaries))
	for i, complementary := range complementaries {
		c := dto.ComplementaryCategoryResponse{
			Id:   complementary.Id,
			Name: complementary.Name,
		}
		complementaryResponse[i] = c
	}
	return complementaryResponse, metadata, nil
}

func (c categoryService) GetComplementaryCategory(id uint64) (*model.Complementary, error) {
	return c.categoryRepository.GetComplementaryCategory(id)
}

func (c categoryService) UpdateCategory(categoryId uint64, req *dto.CategoryRequest) error {
	return c.categoryRepository.UpdateCategory(&model.Category{
		Id:     categoryId,
		Name:   req.Name,
		PicUrl: req.PicUrl,
	})
}

func (c categoryService) UpdateSubCategory(subCategoryId uint64, req *dto.EditSubCategoryRequest) error {
	return c.categoryRepository.UpdateSubCategory(&model.Subcategory{
		Id:   subCategoryId,
		Name: req.Name,
	})
}

func (c categoryService) UpdateComplementaryCategory(complementaryId uint64, req *dto.EditComplementaryCategoryRequest) error {
	return c.categoryRepository.UpdateComplementaryCategory(&model.Complementary{
		Id:   complementaryId,
		Name: req.Name,
	})
}

func (c categoryService) DeleteCategory(categoryId uint64) error {
	return c.categoryRepository.DeleteCategory(&model.Category{
		Id: categoryId,
	})
}

func (c categoryService) DeleteSubCategory(subCategoryId uint64) error {
	return c.categoryRepository.DeleteSubCategory(&model.Subcategory{
		Id: subCategoryId,
	})
}

func (c categoryService) DeleteComplementaryCategory(complementaryId uint64) error {
	return c.categoryRepository.DeleteComplementaryCategory(&model.Complementary{
		Id: complementaryId,
	})
}
