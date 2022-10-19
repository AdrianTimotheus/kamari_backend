package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
)

type ProductRepository interface {
	GetMerchantProduct(merchantId uint64) ([]model.Product, error)
	GetAllProducts(keyword string, filters *helper.Filters) ([]model.ProductList, *helper.Metadata, error)
	GetMerchantProducts(MerchantID uint64, pageOffset int, orderBy string, Limit int, sellerCategoryId int, marketplaceCategoryId int) ([]model.Product, error)
	GetCountMerchantProduct(username string) (int64, error)
	GetProductReviewScore(productID uint64) (count int, score int, err error)
	GetProductReviewScoreDecimal(productId uint64) (int, float32, error)
	GetProductFirstVariant(productID uint64) (*model.ProductVariant, error)
	GetSearchProducts(merchantID uint64, pageOffset int, orderBy string, keyword string, Limit int, sellerCategoryId int, marketplaceCategoryId int) ([]model.Product, error)
	GetTopSellProduct() ([]model.Product, error)
	GetProductDetails(productId uint64) (*model.Product, error)
	GetSimilarProduct(keyword string, category uint64, currentProductId uint64) ([]model.Product, error)
	AddNewProduct(payload model.Product, user model.User) error
	EditProduct(payload model.Product, user model.User) error
	DeleteProduct(id int, user model.User) error
	GetProductPlainDetail(productId uint64) (*model.Product, error)
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
		Where(`to_tsvector('simple', "products"."name") @@ plainto_tsquery('simple', ?) OR ? = ''`, keyword, keyword)

	if filters.CategoryId != nil {
		res = res.Where(`"products"."category_id" = ?`, *filters.CategoryId)
	}
	switch filters.SortBy {
	case "recommended":
		res = res.Order(`"products"."view_count" DESC, "products"."favourite_count" DESC`)
	case "newest":
		res = res.Order(`"products"."created_at" DESC`)
	case "mostbuy":
		res = res.Order(`"products"."unit_sold" DESC`)
	default:
		res = res.Order(`"products"."id" DESC`)
	}

	if err := res.Find(&products).Count(&count).Error; err != nil {
		return nil, nil, err
	}

	res = res.
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&products)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	//err := p.db.
	//	Table("products").
	//	Select("COUNT(*) OVER() AS COUNT").
	//	Limit(1).
	//	Scan(&count).Error
	//if err != nil {
	//	return nil, nil, res.Error
	//}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)

	return products, &metadata, nil
}

func (p productRepository) GetProductReviewScore(productId uint64) (int, int, error) {
	var score int
	var count int
	err := p.db.
		Raw(`SELECT count(product_reviews.id) AS count,COALESCE(ROUND(avg(score)),0) AS score
from product_reviews where order_detail_id in (
    select id
    from order_details
    where order_details.product_variant_id in (
        select id
        from product_variants
        where product_variants.product_id = ?
    )
)`, productId).Row().Scan(&count, &score)
	if err != nil {
		return 0, 0, err
	}

	return count, score, nil
}

func (p productRepository) GetProductReviewScoreDecimal(productId uint64) (int, float32, error) {
	var score float32
	var count int
	err := p.db.
		Raw(`SELECT count(product_reviews.id) AS count,COALESCE(avg(score),0) AS score
from product_reviews where order_detail_id in (
    select id
    from order_details
    where order_details.product_variant_id in (
        select id
        from product_variants
        where product_variants.product_id = ?
    )
)`, productId).Row().Scan(&count, &score)
	if err != nil {
		return 0, 0, err
	}

	return count, score, nil
}

func (p productRepository) GetMerchantProducts(merchantID uint64, pageOffset int, orderBy string, Limit int, sellerCategoryId int, marketplaceCategoryId int) ([]model.Product, error) {
	fmt.Println(pageOffset)
	var products []model.Product
	var err error
	err = nil
	if marketplaceCategoryId != 0 {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Where("category_id", marketplaceCategoryId).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Find(&products).Error
	} else if sellerCategoryId != 0 {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Where("seller_category_id=?", sellerCategoryId).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Find(&products).Error
	} else {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Find(&products).Error
	}

	if err != nil {
		return nil, err
	}
	return products, nil
}

//func (p productRepository) GetCountMerchantProductVariants(username string) (int64, error) {
//	var numOfRecords int64
//	var pvModel model.ProductVariant
//
//	err := p.db.Model(&pvModel).Joins("JOIN users ON users.username=?", username).Count(&numOfRecords).Error
//	if err != nil {
//		return 0, err
//	}
//	return numOfRecords, nil
//}

func (p productRepository) GetCountMerchantProduct(username string) (int64, error) {
	var numOfRecords int64
	var pvModel model.Product

	err := p.db.Model(&pvModel).
		Joins("JOIN users ON users.username=?", username).
		Where("availability=?", true).
		Count(&numOfRecords).Error
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

func (p productRepository) GetSearchProducts(merchantID uint64, pageOffset int, orderBy string, keyword string, Limit int, sellerCategoryId int, marketplaceCategoryId int) ([]model.Product, error) {
	fmt.Println(pageOffset)
	var products []model.Product
	ilikeStatement := fmt.Sprintf("products.name ilike '%%%v%%'", keyword)
	fmt.Println(ilikeStatement)
	var err error
	err = nil
	if marketplaceCategoryId != 0 {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Where(ilikeStatement).
			Where("category_id", marketplaceCategoryId).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Order("unit_sold DESC").
			Find(&products).Error
	} else if sellerCategoryId != 0 {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Where(ilikeStatement).
			Where("seller_category_id", sellerCategoryId).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Order("unit_sold DESC").
			Find(&products).Error
	} else {
		err = p.db.
			Where("merchant_id=?", merchantID).
			Where(ilikeStatement).
			Offset(pageOffset).
			Limit(Limit).
			Order(orderBy).
			Order("unit_sold DESC").
			Find(&products).Error
	}
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

func (p productRepository) GetProductDetails(productId uint64) (*model.Product, error) {
	var product *model.Product
	err := p.db.
		Preload("Merchant").
		Preload("Merchant.User").
		Preload("Category").
		Preload("Subcategory").
		Preload("ProductPhotos").
		Preload("ProductVariants").
		Preload("ProductVariants.ProductPromo", "stock_left > 0").
		Preload("ProductVariants.ProductPromo.PromoType").
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, err
}

func (p productRepository) GetProductPlainDetail(productId uint64) (*model.Product, error) {
	var product *model.Product
	err := p.db.
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return product, err
}

func (p productRepository) GetSimilarProduct(keyword string, categoryId uint64, currentProductId uint64) ([]model.Product, error) {
	var outputModels []model.Product
	//THIS ONLY WORK WITH POSTGRES THAT HAVE pg_trgm EXTENSION
	err := p.db.
		Joins("Merchant").
		Where("similarity(products.name,?)>=0.2", keyword).
		Where("category_id=?", categoryId).
		Where("products.id<>?", currentProductId).
		Limit(24).
		Find(&outputModels).Error
	if err != nil {
		return nil, err
	}
	return outputModels, err
}

func (p productRepository) AddNewProduct(product model.Product, user model.User) error {
	var newProduct model.Product
	var exists bool
	var err error
	var merchant model.Merchant

	if product.ComplementaryCategoryId == 0 {
		err = p.db.Table("categories").
			Joins("join subcategories on categories.id = subcategories.category_id").
			Select("count(*) > 0").
			Where("subcategories.category_id = ?", product.SubcategoryId).
			Where("categories.id = ?", product.CategoryId).
			Find(&exists).
			Error
	} else {
		err = p.db.Table("categories").
			Joins("join subcategories on categories.id = subcategories.category_id").
			Joins("join complementary_categories on subcategories.id = complementary_categories.subcategory_id").
			Select("count(*) > 0").
			Where("complementary_categories.id = ?", product.ComplementaryCategoryId).
			Where("subcategories.id = ?", product.SubcategoryId).
			Where("categories.id = ?", product.CategoryId).
			Find(&exists).
			Error
	}

	if !exists {
		return httperror.BadRequestError("Categories, subcategories, and complementary categories don't match", "CATEGORIES_NOT_MATCH")
	}

	err = p.db.Table("merchants").Where("user_id = ?", user.Id).First(&merchant).Error
	if err != nil {
		return httperror.BadRequestError("User has not registered as a merchant", "USER_NOT_MERCHANT")
	}

	product.MerchantID = merchant.Id

	err = p.db.Table("products").Create(&product).Scan(&newProduct).Error
	if err != nil {
		return err
	}

	return nil
}
func (p productRepository) EditProduct(product model.Product, user model.User) error {
	//var newProduct model.Product
	var exists bool
	var err error
	var merchant model.Merchant

	err = p.db.Table("merchants").Where("user_id = ?", user.Id).First(&merchant).Error
	if err != nil {
		return httperror.BadRequestError("User has not registered as a merchant", "USER_NOT_MERCHANT")
	}

	err = p.db.Table("products").
		Select("count(*) > 0").
		Where("id = ?", product.Id).
		Where("merchant_id = ?", merchant.Id).
		Find(&exists).
		Error

	if !exists {
		return httperror.BadRequestError("Product not found", "PRODUCT_NOT_FOUND")
	}

	if product.ComplementaryCategoryId == 0 {
		err = p.db.Table("categories").
			Joins("join subcategories on categories.id = subcategories.category_id").
			Select("count(*) > 0").
			Where("subcategories.category_id = ?", product.SubcategoryId).
			Where("categories.id = ?", product.CategoryId).
			Find(&exists).
			Error
	} else {
		err = p.db.Table("categories").
			Joins("join subcategories on categories.id = subcategories.category_id").
			Joins("join complementary_categories on subcategories.id = complementary_categories.subcategory_id").
			Select("count(*) > 0").
			Where("complementary_categories.id = ?", product.ComplementaryCategoryId).
			Where("subcategories.id = ?", product.SubcategoryId).
			Where("categories.id = ?", product.CategoryId).
			Find(&exists).
			Error
	}

	if !exists {
		return httperror.BadRequestError("Categories, subcategories, and complementary categories don't match", "CATEGORIES_NOT_MATCH")
	}

	var variant model.ProductVariant
	err = p.db.Table("product_variants").Where("product_id = ?", product.Id).Delete(&variant).Error
	if err != nil {
		return err
	}

	err = p.db.Table("products").Where("id = ?", product.Id).Updates(&product).Error
	if err != nil {
		return err
	}
	return err
}
func (p productRepository) DeleteProduct(id int, user model.User) error {
	var exists bool
	var err error
	var merchant model.Merchant

	err = p.db.Table("merchants").Where("user_id = ?", user.Id).First(&merchant).Error
	if err != nil {
		return httperror.BadRequestError("User has not registered as a merchant", "USER_NOT_MERCHANT")
	}

	err = p.db.Table("products").
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("merchant_id = ?", merchant.Id).
		Where("deleted_at IS NULL").
		Find(&exists).
		Error

	if !exists {
		return httperror.BadRequestError("Product not found", "PRODUCT_NOT_FOUND")
	}

	var newProduct model.Product
	err = p.db.Table("products").Where("id = ?", id).Delete(&newProduct).Error
	if err != nil {
		return err
	}

	var variant model.ProductVariant
	err = p.db.Table("product_variants").Where("product_id = ?", id).Delete(&variant).Error
	if err != nil {
		return err
	}

	return nil
}
