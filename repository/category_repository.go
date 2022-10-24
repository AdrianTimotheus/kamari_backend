package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
)

type CategoryRepository interface {
	SaveCategory(category *model.Category) (*model.Category, error)
	SaveSubCategory(subCategory *model.Subcategory) (*model.Subcategory, error)
	SaveComplementaryCategory(complementary *model.Complementary) (*model.Complementary, error)
	GetCategory(category *model.Category) (*model.Category, error)
	GetCategories(keyword string, filters *helper.Filters) ([]model.Category, *helper.Metadata, error)
	GetSubCategories(keyword string, filters *helper.Filters) ([]model.Subcategory, *helper.Metadata, error)
	GetSubcategory(id uint64) (*model.Subcategory, error)
	GetComplementaryCategories(keyword string, filters *helper.Filters) ([]model.Complementary, *helper.Metadata, error)
	GetComplementaryCategory(id uint64) (*model.Complementary, error)
	UpdateCategory(category *model.Category) error
	UpdateSubCategory(subCategory *model.Subcategory) error
	UpdateComplementaryCategory(complementary *model.Complementary) error
	DeleteCategory(category *model.Category) error
	DeleteSubCategory(subCategory *model.Subcategory) error
	DeleteComplementaryCategory(complementary *model.Complementary) error
}

type categoryRepository struct {
	db *gorm.DB
}

type CaRConfig struct {
	DB *gorm.DB
}

func NewCategoryRepository(c CaRConfig) CategoryRepository {
	return categoryRepository{
		db: c.DB,
	}
}

func (c categoryRepository) isCategoryExists(name string) error {
	var exists bool
	_ = c.db.Table("categories").
		Select("count(*) > 0").
		Where("name = ?", name).
		Find(&exists).Error
	if exists {
		return httperror.BadRequestError("Category name is already exists", "CATEGORY_EXISTS")
	}
	return nil
}

func (c categoryRepository) isSubCategoryExists(name string) error {
	var exists bool
	_ = c.db.Table("subcategories").
		Select("count(*) > 0").
		Where("name = ?", name).
		Find(&exists).Error
	if exists {
		return httperror.BadRequestError("Sub Category name is already exists", "SUB_CATEGORY_EXISTS")
	}
	return nil
}

func (c categoryRepository) isComplementaryCategoryExists(name string) error {
	var exists bool
	_ = c.db.Table("complementary_categories").
		Select("count(*) > 0").
		Where("name = ?", name).
		Find(&exists).Error
	if exists {
		return httperror.BadRequestError("Complementary Category name is already exists", "COMPLEMENTARY_CATEGORY_EXISTS")
	}
	return nil
}

func (c categoryRepository) SaveCategory(category *model.Category) (*model.Category, error) {
	if err := c.isCategoryExists(category.Name); err != nil {
		return nil, err
	}
	result := c.db.
		Create(category)
	return category, result.Error
}

func (c categoryRepository) SaveSubCategory(subCategory *model.Subcategory) (*model.Subcategory, error) {
	if err := c.isSubCategoryExists(subCategory.Name); err != nil {
		return nil, err
	}
	result := c.db.
		Create(subCategory)
	return subCategory, result.Error
}

func (c categoryRepository) SaveComplementaryCategory(complementary *model.Complementary) (*model.Complementary, error) {
	if err := c.isComplementaryCategoryExists(complementary.Name); err != nil {
		return nil, err
	}
	result := c.db.
		Table("complementary_categories").
		Create(complementary)
	return complementary, result.Error
}

func (c categoryRepository) GetCategory(category *model.Category) (*model.Category, error) {
	result := c.db.Find(&category)
	return category, result.Error
}

func (c categoryRepository) GetCategories(keyword string, filters *helper.Filters) ([]model.Category, *helper.Metadata, error) {
	var (
		categories []model.Category
		count      int64
	)

	k := fmt.Sprintf("%%%s%%", keyword)
	condition := c.db.
		Preload("Subcategories").
		Where("name ILIKE ?", k)

	res := condition.
		Order("id ASC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&categories)

	if res.Error != nil {
		return nil, nil, res.Error
	}

	err := condition.
		Select("COUNT(*) OVER() AS COUNT").
		Limit(1).
		Scan(&count).Error
	if err != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return categories, &metadata, nil
}

func (c categoryRepository) GetSubCategories(keyword string, filters *helper.Filters) ([]model.Subcategory, *helper.Metadata, error) {
	var (
		subcategories []model.Subcategory
		count         int64
	)

	k := fmt.Sprintf("%%%s%%", keyword)
	condition := c.db.
		Where("name ILIKE ?", k)

	if filters.CategoryId != nil {
		condition = condition.Where("category_id = ?", filters.CategoryId)
	}

	res := condition.
		Order("id ASC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&subcategories)

	if res.Error != nil {
		return nil, nil, res.Error
	}

	err := condition.
		Select("COUNT(*) OVER() AS COUNT").
		Limit(1).
		Scan(&count).Error
	if err != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return subcategories, &metadata, nil
}

func (c categoryRepository) GetSubcategory(id uint64) (*model.Subcategory, error) {
	var subcategory model.Subcategory
	err := c.db.Where("id = ?", id).First(&subcategory).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

func (c categoryRepository) GetComplementaryCategories(keyword string, filters *helper.Filters) ([]model.Complementary, *helper.Metadata, error) {
	var (
		complementaries []model.Complementary
		count           int64
	)

	k := fmt.Sprintf("%%%s%%", keyword)
	condition := c.db.
		Table("complementary_categories").
		Where("name ILIKE ?", k)

	if filters.SubCategoryId != nil {
		condition = condition.Where("subcategory_id = ?", filters.SubCategoryId)
	}

	res := condition.
		Order("id ASC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&complementaries)

	if res.Error != nil {
		return nil, nil, res.Error
	}

	err := condition.
		Select("COUNT(*) OVER() AS COUNT").
		Limit(1).
		Scan(&count).Error
	if err != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return complementaries, &metadata, nil
}

func (c categoryRepository) GetComplementaryCategory(id uint64) (*model.Complementary, error) {
	var complementary model.Complementary
	err := c.db.Table("complementary_categories").Where("id = ?", id).First(&complementary).Error
	if err != nil {
		return nil, err
	}
	return &complementary, nil
}

func (c categoryRepository) UpdateCategory(category *model.Category) error {
	return c.db.Updates(&category).Error
}

func (c categoryRepository) UpdateSubCategory(subCategory *model.Subcategory) error {
	return c.db.Updates(&subCategory).Error
}

func (c categoryRepository) UpdateComplementaryCategory(complementary *model.Complementary) error {
	return c.db.Table("complementary_categories").Updates(&complementary).Error
}

func (c categoryRepository) DeleteCategory(category *model.Category) error {
	var isUsed bool
	_ = c.db.Table("products").
		Select("count(*) > 0").
		Where("category_id = ? AND deleted_at IS NULL", category.Id).
		Find(&isUsed).Error
	if isUsed {
		return httperror.BadRequestError("Category is used", "CATEGORY_USED")
	}

	var hasChild bool
	_ = c.db.Table("subcategories").
		Select("count(*) > 0").
		Where("category_id = ? AND deleted_at IS NULL", category.Id).
		Find(&hasChild).Error
	if hasChild {
		return httperror.BadRequestError("Category has child", "CATEGORY_HAS_CHILD")
	}

	return c.db.Table("categories").Delete(&category).Error
}

func (c categoryRepository) DeleteSubCategory(subCategory *model.Subcategory) error {
	var isUsed bool
	_ = c.db.Table("products").
		Select("count(*) > 0").
		Where("subcategory_id = ? AND deleted_at IS NULL", subCategory.Id).
		Find(&isUsed).Error
	if isUsed {
		return httperror.BadRequestError("Sub category is used", "SUB_CATEGORY_USED")
	}

	var hasChild bool
	_ = c.db.Table("complementary_categories").
		Select("count(*) > 0").
		Where("subcategory_id = ? AND deleted_at IS NULL", subCategory.Id).
		Find(&hasChild).Error
	if hasChild {
		return httperror.BadRequestError("Sub category has child", "SUB_CATEGORY_HAS_CHILD")
	}

	return c.db.Table("subcategories").Delete(&subCategory).Error
}

func (c categoryRepository) DeleteComplementaryCategory(complementary *model.Complementary) error {
	var isUsed bool
	_ = c.db.Table("products").
		Select("count(*) > 0").
		Where("complementary_category_id = ? AND deleted_at IS NULL", complementary.Id).
		Find(&isUsed).Error
	if isUsed {
		return httperror.BadRequestError("Complementary category is used", "COMPLEMENTARY_CATEGORY_USED")
	}
	return c.db.Table("complementary_categories").Delete(&complementary).Error
}
