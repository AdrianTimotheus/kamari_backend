package repository

import (
	"fmt"
	"gorm.io/gorm"
	"kamari.com/helper"
	"kamari.com/model"
	"time"
)

type ProductPromoRepository interface {
	CreateMerchantPromoTransaction(productPromo *model.ProductPromo, promoItems []model.ProductPromoItem) error
	UpdatePromoTransaction(promo *model.ProductPromo, promoItems []model.ProductPromoItem) error
	EndPromoTransaction(promoId uint64) error
	InsertPromoItems(promoItem []model.ProductPromoItem) error
	GetProductPromos(merchantId uint64, keyword string, filters *helper.Filters) ([]model.ProductPromo, *helper.Metadata, error)
	GetProductPromoDetail(promoId uint64) (*model.ProductPromo, error)
	GetAvailableMerchantProductPromo(merchantId uint64, keyword string, filters *helper.Filters) ([]model.Product, *helper.Metadata, error)
	IsRemovable(promoId uint64) (bool, error)
	DeletePromoItem(itemId uint64) error
	IsValidPromoItem(merchantId, promoId, promoItemId uint64) (bool, error)
}

type productPromoRepository struct {
	db *gorm.DB
}

type PPRConfig struct {
	DB *gorm.DB
}

func NewProductPromoRepository(c PPRConfig) ProductPromoRepository {
	return productPromoRepository{
		db: c.DB,
	}
}

func (p productPromoRepository) CreateMerchantPromoTransaction(productPromo *model.ProductPromo, promoItems []model.ProductPromoItem) error {
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		promoId, err := p.createMerchantPromo(productPromo, tx)
		if err != nil {
			return err
		}
		if err := p.createMerchantPromoItems(promoItems, promoId, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (p productPromoRepository) InsertPromoItems(promoItems []model.ProductPromoItem) error {
	for i := range promoItems {
		if err := p.db.Create(&promoItems[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func (p productPromoRepository) GetProductPromos(merchantId uint64, keyword string, filters *helper.Filters) ([]model.ProductPromo, *helper.Metadata, error) {
	var (
		promos []model.ProductPromo
		count  int64
	)

	k := fmt.Sprintf("%s%%", keyword)
	res := p.db.
		Preload("ProductPromoItems").
		Preload("ProductPromoItems.PromoType").
		Preload("ProductPromoItems.Product").
		Unscoped().
		Where("merchant_id = ?", merchantId).
		Where("title ILIKE ?", k)

	if filters.PromoStatus != nil {
		now := time.Now().Format("2006-01-02")
		switch *filters.PromoStatus {
		case "ongoing":
			res = res.
				Where("start_date <= ?", now).
				Where("end_date >= ?", now)
		case "willCome":
			res = res.
				Where("start_date > ?", now)
		case "hasEnded":
			res = res.
				Where("end_date < ?", now)
		}
	}

	if err := res.Find(&promos).Count(&count).Error; err != nil {
		return nil, nil, err
	}

	res = res.
		Order("updated_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&promos)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)

	return promos, &metadata, nil
}

func (productPromoRepository) createMerchantPromo(promo *model.ProductPromo, tx *gorm.DB) (uint64, error) {
	if err := tx.Create(&promo).Error; err != nil {
		return 0, err
	}

	return promo.Id, nil
}

func (productPromoRepository) createMerchantPromoItems(promoItems []model.ProductPromoItem, promoId uint64, tx *gorm.DB) error {
	for i := range promoItems {
		promoItems[i].ProductPromoId = promoId
		if err := tx.Create(&promoItems[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func (p productPromoRepository) GetProductPromoDetail(promoId uint64) (*model.ProductPromo, error) {
	var promo model.ProductPromo

	err := p.db.
		Preload("ProductPromoItems").
		Preload("ProductPromoItems.PromoType").
		Preload("ProductPromoItems.Product").
		Preload("ProductPromoItems.Product.ProductVariants").
		Unscoped().
		Where("id = ?", promoId).
		First(&promo).Error
	if err != nil {
		return nil, err
	}

	return &promo, nil
}

func (p productPromoRepository) IsRemovable(promoId uint64) (bool, error) {
	now := time.Now().Format("2006-01-02")
	res := p.db.
		Where("id = ?", promoId).
		Where("start_date > ?", now).
		Limit(1).
		Find(&model.ProductPromo{})
	if res.Error != nil {
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (p productPromoRepository) DeletePromoItem(itemId uint64) error {
	return p.db.
		Where("id = ?", itemId).
		Delete(&model.ProductPromoItem{}).Error
}

func (p productPromoRepository) IsValidPromoItem(merchantId, promoId, promoItemId uint64) (bool, error) {
	var count int64

	res := p.db.
		Table("product_promo_items").
		Select("COUNT(*)").
		Joins(`LEFT JOIN product_promos ON "product_promo_items"."product_promo_id" = "product_promos"."id"`).
		Where(`"product_promo_items"."id" = ?`, promoItemId).
		Where(`"product_promo_items"."deleted_at" IS NULL`).
		Where(`"product_promos"."merchant_id" = ?`, merchantId).
		Where(`"product_promos"."id" = ?`, promoId).
		Limit(1).
		Scan(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (p productPromoRepository) UpdatePromoTransaction(promo *model.ProductPromo, promoItems []model.ProductPromoItem) error {
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := p.updateMerchantPromo(promo, tx); err != nil {
			return err
		}
		if err := p.updateMerchantPromoItems(promoItems, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (p productPromoRepository) GetAvailableMerchantProductPromo(merchantId uint64, keyword string, filters *helper.Filters) ([]model.Product, *helper.Metadata, error) {
	var (
		count      int64
		productIds []int
		promoIds   []int
		products   []model.Product
	)

	k := fmt.Sprintf("%s%%", keyword)

	if err := p.db.
		Table("products").
		Select("id").
		Where("merchant_id = ?", merchantId).
		Scan(&productIds).
		Error; err != nil {
		return nil, nil, err
	}

	if err := p.db.
		Table("product_promo_items").
		Select("product_id").
		Where("product_id IN ?", productIds).
		Where("deleted_at IS NULL").
		Scan(&promoIds).
		Error; err != nil {
		return nil, nil, err
	}

	c := p.db.Preload("ProductVariants")

	var res *gorm.DB
	if len(promoIds) < 1 {
		if res = c.Where("id in ?", productIds).
			Find(&products).Count(&count); res.Error != nil {
			return nil, nil, res.Error
		}
	} else {
		if res = c.
			Where("id NOT IN ?", promoIds).
			Where("merchant_id = ?", merchantId).
			Where("name ILIKE ?", k).
			Find(&products).Count(&count); res.Error != nil {
			return nil, nil, res.Error
		}
	}

	res = res.
		Order("updated_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&products)
	if res.Error != nil {
		return nil, nil, res.Error
	}

	metadata := helper.CalculateMetaData(int(count), filters.Page, filters.Limit)

	return products, &metadata, nil
}

func (p productPromoRepository) EndPromoTransaction(promoId uint64) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		promoItems, err := p.getPromoItems(promoId)
		if err != nil {
			return err
		}

		for _, item := range promoItems {
			if err = p.deletePromoItem(tx, item.Id); err != nil {
				return err
			}
		}

		if err = tx.Delete(&model.ProductPromo{Id: promoId}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (productPromoRepository) deletePromoItem(tx *gorm.DB, itemId uint64) error {
	return tx.
		Where("id = ?", itemId).
		Delete(&model.ProductPromoItem{}).Error
}

func (p productPromoRepository) getPromoItems(promoId uint64) ([]model.ProductPromoItem, error) {
	var promoItems []model.ProductPromoItem

	err := p.db.
		Where("promo_id = ?", promoId).
		Find(&promoItems).Error
	if err != nil {
		return nil, err
	}

	return promoItems, nil
}

func (productPromoRepository) updateMerchantPromo(promo *model.ProductPromo, tx *gorm.DB) error {
	if err := tx.Updates(&promo).Error; err != nil {
		return err
	}

	return nil
}

func (productPromoRepository) updateMerchantPromoItems(promoItems []model.ProductPromoItem, tx *gorm.DB) error {
	for _, item := range promoItems {
		if err := tx.Updates(&item).Error; err != nil {
			return err
		}
	}

	return nil
}
