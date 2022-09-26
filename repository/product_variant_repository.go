package repository

import (
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProductVariant interface {
	GetProductVariantByID(id int) (model.ProductVariants, error)
}

type productVariant struct {
	db *gorm.DB
}

type PVConfig struct {
}
