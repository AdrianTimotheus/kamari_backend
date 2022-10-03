package repository

import (
	"gorm.io/gorm"
	"kamari.com/dto"
)

type CartRepository interface {
	GetCartHoverItems(userId uint64) (*dto.CartHoverRes, error)
}

type cartRepository struct {
	db *gorm.DB
}

type CRConfig struct {
	DB *gorm.DB
}

func NewCartRepository(c CRConfig) CartRepository {
	return cartRepository{db: c.DB}
}

func (c cartRepository) GetCartHoverItems(userId uint64) (*dto.CartHoverRes, error) {
	var cartHoverRes dto.CartHoverRes
	var count uint64
	var items []dto.CartHoverItem

	result := c.db.Table("product_variants").
		Select("products.name", "product_variants.price", "product_variants.pic_url").
		Joins("JOIN products on product_variants.product_id = products.id").
		Joins("JOIN carts on carts.variant_id = product_variants.id").
		Where("carts.user_id = ?", userId).
		Order("carts.updated_at desc").
		Limit(5).
		Scan(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	result = c.db.Table("carts").
		Select("count(*)").
		Where("user_id = ?", userId).
		Scan(&count)

	if result.Error != nil {
		return nil, result.Error
	}

	if count < 5 {
		cartHoverRes.RemainingItem = 0
	} else {
		cartHoverRes.RemainingItem = count - 5
	}

	cartHoverRes.Items = items

	return &cartHoverRes, nil
}
