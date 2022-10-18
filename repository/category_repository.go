package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/model"
)

type CategoryRepository interface {
	GetCategory(category *model.Category) (*model.Category, error)
	GetCategories(keyword string, filters *helper.Filters) ([]model.Category, *helper.Metadata, error)
	GetSubcategories() ([]model.Subcategory, error)
	GetSubcategory(id uint64) (*model.Subcategory, error)
	GetComplementaries() ([]model.Complementary, error)
	GetComplementary(id uint64) (*model.Complementary, error)
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

func (c categoryRepository) GetSubcategories() ([]model.Subcategory, error) {
	var subcategories []model.Subcategory
	err := c.db.Find(&subcategories).Error
	if err != nil {
		return nil, err
	}
	return subcategories, nil
}

func (c categoryRepository) GetSubcategory(id uint64) (*model.Subcategory, error) {
	var subcategory model.Subcategory
	err := c.db.Where("id = ?", id).First(&subcategory).Error
	if err != nil {
		return nil, err
	}
	return &subcategory, nil
}

func (c categoryRepository) GetComplementaries() ([]model.Complementary, error) {
	var complementaries []model.Complementary
	err := c.db.Table("complementary_categories").Find(&complementaries).Error
	if err != nil {
		return nil, err
	}
	return complementaries, nil
}

func (c categoryRepository) GetComplementary(id uint64) (*model.Complementary, error) {
	var complementary model.Complementary
	err := c.db.Table("complementary_categories").Where("id = ?", id).First(&complementary).Error
	if err != nil {
		return nil, err
	}
	return &complementary, nil
}
