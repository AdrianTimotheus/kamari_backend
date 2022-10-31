package repository

import (
	"fmt"
	"math"
	"net/url"
	"strconv"

	"gorm.io/gorm"
	"kamari.com/dto"
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
	GetSellerCategories(userId uint64, params url.Values) ([]dto.GetSellerCategoriesRes, uint64, error)
	GetSellerPageCategories(username string) ([]model.SellerCategory, error)
	GetSellerCategory(userId uint64, sellerCategoryId uint64) (string, []dto.SellerCategoriesItemsRes, error)
	UpdateSellerCategory(userId uint64, category model.SellerCategory) error
	CreateSellerCategory(userId uint64, category model.SellerCategory) error
	DeleteSellerCategory(userId uint64, sellerCategoryId uint64) error
	GetAvailableProducts(userId uint64, params url.Values) ([]dto.SellerCategoriesAvailableItemsRes, uint64, error)
	AddProductToSellerCategory(userId uint64, req dto.AddRemoveProductFromCategory, sellerCategoryId uint64) error
	DeleteProductFromSellerCategory(userId uint64, req dto.AddRemoveProductFromCategory) error
	RemoveSpecificProductFromSellerCategory(userId uint64, productId uint64) error
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

func (c categoryRepository) GetSellerCategories(userId uint64, params url.Values) ([]dto.GetSellerCategoriesRes, uint64, error) {
	var productsCountBySellerCategories []dto.GetSellerCategoriesRes
	var merchant model.Merchant
	var count int64

	keywordParams := params["keyword"]
	keyword := ""
	if keywordParams != nil {
		keyword = "%" + keywordParams[0] + "%"
	}

	limit, _ := strconv.Atoi(params["limit"][0])
	page, _ := strconv.Atoi(params["page"][0])
	offset := limit * (page - 1)

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return nil, 0, err
	}

	if keyword == "" {
		err = c.db.Table("seller_categories").
			Select("seller_categories.id, seller_categories.name, count(products.id)").
			Joins("left join products on seller_categories.id = products.seller_category_id").
			Where("seller_categories.merchant_id = ?", merchant.Id).
			Where("seller_categories.deleted_at IS NULL").
			Group("seller_categories.name, seller_categories.id").
			Count(&count).Error
		err = c.db.Table("seller_categories").
			Select("seller_categories.id, seller_categories.name, count(products.id)").
			Joins("left join products on seller_categories.id = products.seller_category_id").
			Where("seller_categories.merchant_id = ?", merchant.Id).
			Where("seller_categories.deleted_at IS NULL").
			Group("seller_categories.name, seller_categories.id").
			Limit(limit).
			Offset(offset).
			Find(&productsCountBySellerCategories).
			Error
	} else {
		//SELECT seller_categories.name, seller_category_id, count(seller_categories.name) from seller_categories join products on seller_categories.id = products.seller_category_id group by seller_category_id, seller_categories.name;
		err = c.db.Table("seller_categories").
			Select("seller_categories.id, seller_categories.name, count(products.id)").
			Joins("left join products on seller_categories.id = products.seller_category_id").
			Where("seller_categories.merchant_id = ?", merchant.Id).
			Where("seller_categories.name ILIKE ?", keyword).
			Where("seller_categories.deleted_at IS NULL").
			Group("seller_categories.name, seller_categories.id").
			Count(&count).Error
		fmt.Println(count)

		err = c.db.Table("seller_categories").
			Select("seller_categories.id, seller_categories.name, count(products.id)").
			Joins("left join products on seller_categories.id = products.seller_category_id").
			Where("seller_categories.merchant_id = ?", merchant.Id).
			Where("seller_categories.name ILIKE ?", keyword).
			Where("seller_categories.deleted_at IS NULL").
			Group("seller_categories.name, seller_categories.id").
			Limit(limit).
			Offset(offset).
			Find(&productsCountBySellerCategories).
			Error
	}

	if err != nil {
		return nil, 0, err
	}

	//pp, _ := json.MarshalIndent(productsCountBySellerCategories, "", "\t")
	//fmt.Println("testtt", string(pp))
	return productsCountBySellerCategories, uint64(math.Ceil(float64(count) / float64(limit))), nil
}

func (c categoryRepository) GetSellerPageCategories(username string) ([]model.SellerCategory, error) {
	var categories []model.SellerCategory
	result := c.db.
		Table("seller_categories").
		Joins("JOIN merchants on seller_categories.merchant_id = merchants.id").
		Joins("JOIN users on merchants.user_id = users.id").
		Where("users.username = ?", username).
		Find(&categories)
	return categories, result.Error
}

func (c categoryRepository) GetSellerCategory(userId uint64, sellerCategoryId uint64) (string, []dto.SellerCategoriesItemsRes, error) {
	var products []dto.SellerCategoriesItemsRes
	var sellerCategory model.SellerCategory
	var merchant model.Merchant
	var exists bool

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return "", nil, err
	}

	err = c.db.Table("seller_categories").
		Select("count(*) > 0").
		Where("merchant_id = ?", merchant.Id).
		Where("id = ?", sellerCategoryId).
		Find(&exists).Error

	if !exists {
		return "", nil, httperror.BadRequestError("Category not found", "CATEGORY_NOT_FOUND")
	}

	err = c.db.Table("products").
		//Select("products.id, name, description, thumbnail_url, seller_category_id").
		//Joins("JOIN product_variants ON product_variants.product_id = products.id ").
		Preload("Variants").
		Where("seller_category_id = ?", sellerCategoryId).
		Where("merchant_id = ?", merchant.Id).
		Where("deleted_at IS NULL").
		Find(&products).
		Error

	if err != nil {
		return "", nil, err
	}

	err = c.db.Table("seller_categories").Where("id = ?", sellerCategoryId).First(&sellerCategory).Error

	if err != nil {
		return "", nil, err
	}

	return sellerCategory.Name, products, nil
}

func (c categoryRepository) UpdateSellerCategory(userId uint64, sellerCategory model.SellerCategory) error {
	var exists bool
	var merchant model.Merchant

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return err
	}

	err = c.db.Table("seller_categories").
		Select("count(*) > 0").
		Where("merchant_id = ?", merchant.Id).
		Where("id = ?", sellerCategory.Id).
		Find(&exists).Error

	if !exists {
		return httperror.BadRequestError("Category not found", "CATEGORY_NOT_FOUND")
	}

	err = c.db.Table("seller_categories").Where("id = ?", sellerCategory.Id).Update("name", sellerCategory.Name).Error
	if err != nil {
		return err
	}
	return nil
}

func (c categoryRepository) CreateSellerCategory(userId uint64, category model.SellerCategory) error {
	var merchant model.Merchant

	err := c.db.Table("merchants").Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return err
	}

	category.MerchantId = merchant.Id

	err = c.db.Table("seller_categories").
		Create(&category).
		Error

	if err != nil {
		return err
	}
	return nil
}

func (c categoryRepository) DeleteSellerCategory(userId uint64, sellerCategoryId uint64) error {
	//if category contains product prevent delete
	var exists bool
	var merchant model.Merchant
	var sellerCategory model.SellerCategory

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return err
	}

	err = c.db.Table("seller_categories").
		Select("count(*) > 0").
		Where("merchant_id = ?", merchant.Id).
		Where("id = ?", sellerCategoryId).
		Find(&exists).Error

	if !exists {
		return httperror.BadRequestError("Category not found", "CATEGORY_NOT_FOUND")
	}

	err = c.db.Table("products").
		Select("count(*) > 0").
		Where("seller_category_id = ?", sellerCategoryId).
		Find(&exists).Error

	if exists {
		return httperror.BadRequestError("Cannot delete used category", "CATEGORY_USED_ERROR")
	}

	sellerCategory.Id = sellerCategoryId
	err = c.db.Table("seller_categories").Where("id = ?", sellerCategoryId).Delete(&sellerCategory).Error
	if err != nil {
		return err
	}
	return nil
}

func (c categoryRepository) GetAvailableProducts(userId uint64, params url.Values) ([]dto.SellerCategoriesAvailableItemsRes, uint64, error) {
	var products []dto.SellerCategoriesAvailableItemsRes
	var merchant model.Merchant
	var count int64

	keywordParams := params["keyword"]
	keyword := ""
	if keywordParams != nil {
		keyword = "%" + keywordParams[0] + "%"
	}

	limit, _ := strconv.Atoi(params["limit"][0])
	page, _ := strconv.Atoi(params["page"][0])
	offset := limit * (page - 1)

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		return nil, 0, err
	}

	if keyword == "" {
		err = c.db.Table("products").
			Preload("Variants").
			Where("merchant_id = ?", merchant.Id).
			Where("seller_category_id IS NULL").
			Find(&products).
			Count(&count).
			Error
		err = c.db.Table("products").
			Preload("Variants").
			Where("merchant_id = ?", merchant.Id).
			Where("seller_category_id IS NULL").
			Limit(limit).
			Offset(offset).
			Find(&products).
			Error
	} else {
		//SELECT seller_categories.name, seller_category_id, count(seller_categories.name) from seller_categories join products on seller_categories.id = products.seller_category_id group by seller_category_id, seller_categories.name;
		err = c.db.Table("products").
			Preload("Variants").
			Where("merchant_id = ?", merchant.Id).
			Where("seller_category_id IS NULL").
			Where("name ILIKE ?", keyword).
			Find(&products).
			Count(&count).
			Error

		err = c.db.Table("products").
			Preload("Variants").
			Where("merchant_id = ?", merchant.Id).
			Where("seller_category_id IS NULL").
			Where("name ILIKE ?", keyword).
			Limit(limit).
			Offset(offset).
			Find(&products).
			Error
	}

	if err != nil {
		return nil, 0, err
	}

	return products, uint64(math.Ceil(float64(count) / float64(limit))), nil
}

func (c categoryRepository) AddProductToSellerCategory(userId uint64, req dto.AddRemoveProductFromCategory, sellerCategoryId uint64) error {
	var merchant model.Merchant
	var exists bool
	var filteredId []uint64

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error

	err = c.db.Table("seller_categories").
		Select("count(*) > 0").
		Where("merchant_id = ?", merchant.Id).
		Where("id = ?", sellerCategoryId).
		Find(&exists).Error

	if !exists {
		return httperror.BadRequestError("Seller category not found", "SELLER_CATEGORY_NOT_FOUND")
	}

	for _, id := range req.ItemsId {
		err = c.db.Table("products").
			Select("count(*) > 0").
			Where("merchant_id = ?", merchant.Id).
			Where("id = ?", id).
			Find(&exists).Error

		if exists {
			filteredId = append(filteredId, id)
		}
	}

	if len(filteredId) == 0 {
		return httperror.BadRequestError("No product belongs to current seller", "NO_VALID_PRODUCTS")
	}

	err = c.db.Table("products").Where("id IN ?", filteredId).Update("seller_category_id", sellerCategoryId).Error
	if err != nil {
		return err
	}
	return nil
}

func (c categoryRepository) DeleteProductFromSellerCategory(userId uint64, req dto.AddRemoveProductFromCategory) error {
	var merchant model.Merchant
	var exists bool
	var filteredId []uint64

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error

	for _, id := range req.ItemsId {
		err = c.db.Table("products").
			Select("count(*) > 0").
			Where("merchant_id = ?", merchant.Id).
			Where("id = ?", id).
			Find(&exists).Error

		if exists {
			filteredId = append(filteredId, id)
		}
	}

	err = c.db.Table("products").Where("id IN ?", filteredId).Update("seller_category_id", gorm.Expr("NULL")).Error
	if err != nil {
		return err
	}
	return nil
}

func (c categoryRepository) RemoveSpecificProductFromSellerCategory(userId uint64, productId uint64) error {
	var merchant model.Merchant
	var exists bool

	err := c.db.Where("user_id = ?", userId).First(&merchant).Error

	err = c.db.Table("products").
		Select("count(*) > 0").
		Where("merchant_id = ?", merchant.Id).
		Where("id = ?", productId).
		Find(&exists).Error

	if !exists {
		return httperror.BadRequestError("Product not found", "PRODUCT_NOT_FOUND")
	}

	err = c.db.Table("products").Where("id = ?", productId).Update("seller_category_id", gorm.Expr("NULL")).Error
	if err != nil {
		return err
	}
	return nil
}
