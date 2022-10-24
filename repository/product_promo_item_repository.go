package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/httperror"
	"kamari.com/model"
	"time"
)

type ProductPromoItemRepository interface {
	GetProductPromoItemByProductId(productId uint64) (*model.ProductPromoItem, error)
}

type productPromoItemRepository struct {
	db *gorm.DB
}

type PPIRConfig struct {
	DB *gorm.DB
}

func NewProductPromoItemRepository(c PPIRConfig) ProductPromoItemRepository {
	return &productPromoItemRepository{db: c.DB}
}

func (p productPromoItemRepository) GetProductPromoItemByProductId(productId uint64) (*model.ProductPromoItem, error) {
	var pPIModel model.ProductPromoItem
	now := time.Now().Format("2006-01-02")
	fmt.Println(now)
	err := p.db.
		Joins("ProductPromo").
		Where("start_date  <= ?", now).
		Where("end_date >= ?", now).
		Where("product_promo_items.stock_left > 0").
		Where("product_promo_items.product_id=?", productId).
		First(&pPIModel).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
	}
	return &pPIModel, nil
}
