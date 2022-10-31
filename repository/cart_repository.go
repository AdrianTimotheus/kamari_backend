package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
)

type CartRepository interface {
	GetCartHoverItems(userId uint64) (*dto.CartHoverRes, error)
	GetCartItems(userId uint64) ([]dto.CartItemRaw, error)
	CreateCartItem(userId uint64, req dto.AddCartItemRequest) error
	UpdateCartItem(id int, userId uint64, req dto.UpdateCartItemRequest) error
	DeleteCartItem(id int, userId uint64) error
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

func (c cartRepository) GetCartItems(userId uint64) ([]dto.CartItemRaw, error) {
	var rawItems []dto.CartItemRaw

	err := c.db.Raw("select products.name ,\n"+
		"product_variants.main_variant as main_variant, \n"+
		"product_variants.main_variant_item as main_variant_item,\n"+
		"product_variants.sub_variant as sub_variant,\n"+
		"product_variants.sub_variant_item as sub_variant_item,\n"+
		"carts.quantity ,\n"+
		"products.name as name,\n"+
		"product_variants.pic_url as pic_url ,\n"+
		"product_variants.price as price,\n"+
		"merchants.id as merchant_id, \n"+
		"merchants.name as merchant_name, \n"+
		"product_variants.id as variant_id, \n"+
		"carts.id as id, \n"+
		"merchants.pic_url as merchant_pic_url, \n"+
		"product_promos.start_date as start_date, \n"+
		"product_promos.end_date as end_date, \n"+
		"product_promo_items.stock_left as stock_left, \n"+
		"product_promo_items.maximum_item as maximum_item, \n"+
		"product_promo_items.promo_type_id as promo_type_id, \n"+
		"product_promo_items.discount_amount as discount_amount,	 \n"+
		"products.minimum_order as minimum_order, \n"+
		"products.maximum_order as maximum_order \n"+
		"from carts inner join product_variants on carts.variant_id = product_variants.id \n"+
		"inner join products on products.id = product_variants.product_id  \n"+
		"inner join merchants on merchants.id = products.merchant_id \n"+
		"left join product_promo_items on product_promo_items.product_id = products.id \n"+
		"left join product_promos on product_promo_items.product_promo_id = product_promos.id \n"+
		"where carts.user_id = ?", userId).Scan(&rawItems).Error

	//var promo model.ProductPromoItem
	//
	//err := c.db.Table("product_promos").Where("")

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
		default:
			return nil, httperror.InternalServerError(err.Error())
		}
		//return nil, httperror.InternalServerError("There was a problem getting cart items")
	}

	return rawItems, nil
}

func (c cartRepository) CreateCartItem(userId uint64, req dto.AddCartItemRequest) error {
	var exists bool
	//var cartItem model.CartItem
	var err error
	var productVariant model.ProductVariant

	c.db.Table("product_variants").Where("id = ?", req.VariantId).First(&productVariant)
	fmt.Println(productVariant.Stock)
	fmt.Println(req.Quantity)
	if uint64(productVariant.Stock) < req.Quantity {
		return httperror.BadRequestError("Quantity exceeds stock", "QUANTITY_EXCEEDS_STOCK")
	}

	c.db.Table("carts").
		Select("count(*) > 0").
		Where("user_id = ?", userId).
		Where("variant_id = ?", req.VariantId).
		Find(&exists)

	if exists {
		err = c.db.Table("carts").
			Where("user_id = ?", userId).
			Where("variant_id = ?", req.VariantId).
			Update("quantity", gorm.Expr("quantity + ?", req.Quantity)).Error

		if err != nil {
			return httperror.InternalServerError("Error updating quantity")
		}
		return nil
	}

	err = c.db.Table("carts").Create(&model.CartItem{
		UserId:    userId,
		VariantId: req.VariantId,
		Quantity:  req.Quantity,
	}).Error

	if err != nil {
		return httperror.BadRequestError("Variant not found", "VARIANT_NOT_FOUND")
	}

	return nil
}

func (c cartRepository) UpdateCartItem(id int, userId uint64, req dto.UpdateCartItemRequest) error {
	var exists bool
	//var cartItem model.CartItem
	var err error

	c.db.Table("carts").
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Find(&exists)

	if exists {
		err = c.db.Table("carts").
			Where("id = ?", id).
			Update("quantity", req.Quantity).Error

		if err != nil {
			return httperror.InternalServerError("Error updating quantity")
		}
		return nil

	}

	return httperror.BadRequestError("Cart item not found", "CART_ITEM_NOT_FOUND")
}

func (c cartRepository) DeleteCartItem(id int, userId uint64) error {
	var exists bool
	var cartItem model.CartItem
	var err error

	c.db.Table("carts").
		Select("count(*) > 0").
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Find(&exists)

	if exists {
		err = c.db.Unscoped().Table("carts").Where("id = ?", id).Delete(&cartItem).Error

		if err != nil {
			return httperror.InternalServerError("Failed to remove item from cart")
		}
		return nil
	}

	return httperror.BadRequestError("Cart item not found", "CART_ITEM_NOT_FOUND")
}
