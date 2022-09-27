package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/model"
)

type ProductRepository interface {
	GetMerchantProduct(merchantId uint64) ([]model.Product, error)
	GetAllProducts(keyword string, filters *helper.Filters) ([]model.ProductList, *helper.Metadata, error)
	GetMerchantProducts(MerchantID uint64, pageOffset int, orderBy string) ([]model.Product, error)
	GetCountMerchantProducts(username string) (int64, error)
	GetProductReviewScore(productID uint64) (int, error)
	GetProductFirstVariant(productID uint64) (*model.ProductVariant, error)
	GetSearchProducts(merchantID uint64, pageOffset int, orderBy string, keyword string) ([]model.Product, error)
	GetTopSellProduct() ([]model.Product, error)
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

func (p productRepository) GetMerchantProduct(merchantId uint64) ([]model.Product, error) {
	var product []model.Product
	err := p.db.
		Where("merchant_id = ?", merchantId).
		Find(&product).Error

	if err != nil {
		return product, err
	}
	return product, nil
}

func (p productRepository) GetAllProducts(keyword string, filters *helper.Filters) ([]model.ProductList, *helper.Metadata, error) {
	var (
		products []model.ProductList
		count    int64
	)

	k := fmt.Sprintf("%%%s%%", keyword)
	res := p.db.
		Select(
			"cities.name AS city",
			"products.id",
			"products.name",
			"products.favourite_count",
			"products.thumbnail_url",
			"products.unit_sold",
			"products.category_id",
			"products.availability",
			"products.merchant_id",
			"users.username AS seller",
		).
		Joins("Merchant").
		Joins("Category").
		Joins(`JOIN users ON users.id = "Merchant"."user_id"`).
		Joins(`LEFT JOIN user_addresses ON user_addresses.id = "Merchant"."user_address_id"`).
		Joins(`LEFT JOIN cities ON cities.id = user_addresses.city_id`).
		Where(`"products"."name" ILIKE ?`, k)

	if filters.CategoryId != nil {
		res = res.Where(`"products"."category_id" = ?`, filters.CategoryId)
	}
	switch filters.SortBy {
	case "recommended":
		res = res.Order(`"products"."view_count" DESC, "products"."favourite_count" DESC`)
	case "newest":
		res = res.Order(`"products"."created_at" DESC`)
	case "mostBuy":
		res = res.Order(`"products"."unit_sold" DESC`)
	default:
		res = res.Order(`"products"."id" DESC`)
	}

	res = res.
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&products)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	res = res.Count(&count)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)

	return products, &metadata, nil
}

func (p productRepository) GetProductReviewScore(productId uint64) (int, error) {
	var score int
	err := p.db.
		Raw(`SELECT COALESCE(ROUND(avg(score)),0) AS score 
				from product_reviews where order_detail_id in (
					select id 
					from order_details
					where order_details.product_variant_id in (
						select id 
						from product_variants
						where product_variants.product_id = ?
					)
				)`, productId).Scan(&score).Error
	if err != nil {
		return 0, err
	}

	return score, nil
}
func (p productRepository) GetMerchantProducts(merchantID uint64, pageOffset int, orderBy string) ([]model.Product, error) {
	fmt.Println(pageOffset)
	var products []model.Product
	err := p.db.
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

func (p productRepository) GetCountMerchantProducts(username string) (int64, error) {
	var numOfRecords int64
	var pvModel model.ProductVariant

	err := p.db.Model(&pvModel).Joins("JOIN users ON users.username=?", username).Count(&numOfRecords).Error
	if err != nil {
		return 0, err
	}
	return numOfRecords, nil
}

func (p productRepository) GetProductFirstVariant(productID uint64) (*model.ProductVariant, error) {
	var pvModel model.ProductVariant
	err := p.db.Where("product_id=?", productID).First(&pvModel).Error
	if err != nil {
		return nil, err
	}
	return &pvModel, nil
}

func (p productRepository) GetSearchProducts(merchantID uint64, pageOffset int, orderBy string, keyword string) ([]model.Product, error) {
	fmt.Println(pageOffset)
	var products []model.Product
	ilikeStatement := fmt.Sprintf("products.name ilike '%%%v%%'", keyword)
	err := p.db.
		Where("merchant_id=?", merchantID).
		Where(ilikeStatement).
		Offset(pageOffset).
		Limit(20).
		Order(orderBy).
		Order("unit_sold DESC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p productRepository) GetTopSellProduct() ([]model.Product, error) {
	var models []model.Product
	err := p.db.
		Joins("Merchant").
		Order("unit_sold DESC").
		Limit(20).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
