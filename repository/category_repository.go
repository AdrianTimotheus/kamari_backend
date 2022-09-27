package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/model"
)

type CategoryRepository interface {
	GetCategories(keyword string, filters *helper.Filters) ([]model.Category, *helper.Metadata, error)
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

func (c categoryRepository) GetCategories(keyword string, filters *helper.Filters) ([]model.Category, *helper.Metadata, error) {
	var (
		categories []model.Category
		count      int64
	)

	k := fmt.Sprintf("%%%s%%", keyword)
	res := c.db.
		Preload("Subcategories").
		Where("name ILIKE ?", k).
		Order("id ASC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&categories)

	if res.Error != nil {
		return nil, nil, res.Error
	}

	res = res.Count(&count)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)
	return categories, &metadata, nil
}
