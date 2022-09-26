package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/model"
)

type ProductRepository interface {
	GetMerchantProduct(merchantId uint64) ([]model.Product, error)
	GetMerchantProducts(MerchantID uint64, pageOffset int, orderBy string) ([]model.Product, error)
	GetCountMerchantProducts(username string) (int64, error)
	GetProductReviewScore(productID uint64) (int, error)
	GetProductFirstVariant(productID uint64) (*model.ProductVariants, error)
}

type productRepository struct {
	db *gorm.DB
}

type ProdConfig struct {
	DB *gorm.DB
}

func NewProductRepository(c ProdConfig) ProductRepository {
	return productRepository{
		db: c.DB,
	}
}

func (c productRepository) GetMerchantProduct(merchantId uint64) ([]model.Product, error) {
	var product []model.Product
	err := c.db.
		Where("merchant_id = ?", merchantId).
		Find(&product).Error

	if err != nil {
		return product, err
	}
	return product, nil
}

func (c productRepository) GetMerchantProducts(merchantID uint64, pageOffset int, orderBy string) ([]model.Product, error) {
	fmt.Println(pageOffset)
	var products []model.Product
	err := c.db.
		Where("merchant_id=?", merchantID).
		Offset(pageOffset).
		Limit(20).
		Order(orderBy).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (c productRepository) GetProductReviewScore(productID uint64) (int, error) {
	var score int
	err := c.db.Raw("select COALESCE(ROUND(avg(score)),0) AS score from product_reviews where order_detail_id  in (select id from order_details where order_details.product_variant_id in (select id from product_variants where product_variants.product_id=1))").Scan(&score).Error
	if err != nil {
		fmt.Println("masuk")
		return 0, err
	}
	return score, nil
}

func (c productRepository) GetCountMerchantProducts(username string) (int64, error) {
	var numOfRecords int64
	var pvModel model.ProductVariants

	err := c.db.Model(&pvModel).Joins("JOIN users ON users.username=?", username).Count(&numOfRecords).Error
	if err != nil {
		return 0, err
	}
	return numOfRecords, nil
}

func (c productRepository) GetProductFirstVariant(productID uint64) (*model.ProductVariants, error) {
	var pvModel model.ProductVariants
	err := c.db.Where("product_id=?", productID).First(&pvModel).Error
	if err != nil {
		return nil, err
	}
	return &pvModel, nil
}
