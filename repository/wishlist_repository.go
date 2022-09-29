package repository

import (
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/model"
)

type WishlistRepository interface {
	Favorite(userId uint64, productId uint64) (*dto.StatusMessageResponse, error)
}

type wishlistRepository struct {
	db *gorm.DB
}

type WLConfig struct {
	DB *gorm.DB
}

func NewWishlistRepository(c WLConfig) WishlistRepository {
	return wishlistRepository{db: c.DB}
}

func (w wishlistRepository) Favorite(userId uint64, productId uint64) (*dto.StatusMessageResponse, error) {
	var exists bool
	var wishlist model.Wishlist
	var product model.Product

	w.db.Model(&wishlist).
		Select("count(*) > 0").
		Where("user_id = ?", userId).
		Where("product_id = ?", productId).
		Find(&exists)

	wishlist = model.Wishlist{
		UserId:    userId,
		ProductId: productId,
	}

	err := w.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if !exists {

			if err := tx.Model(&product).Where("id = ?", productId).Update("favourite_count", gorm.Expr("favourite_count + ?", 1)).Error; err != nil {
				return err
			}

			if err := tx.Create(&wishlist).Error; err != nil {
				return err
			}
			return nil
		}

		if err := tx.Unscoped().Where("user_id = ?", userId).Where("product_id = ?", productId).Delete(&wishlist).Error; err != nil {
			return err
		}

		if err := tx.Model(&product).Where("id = ?", productId).Update("favourite_count", gorm.Expr("favourite_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	//if exists {
	//	err := w.db.Unscoped().Where("user_id = ?", userId).Where("product_id = ?", productId).Delete(&wishlist).Error
	//	if err != nil {
	//		return nil, httperror.InternalServerError("Failed to delete")
	//	}
	//	return &dto.StatusMessageResponse{Status: !exists}, nil
	//}
	//err = w.db.Create(&wishlist).Error
	//if err != nil {
	//	return nil, httperror.InternalServerError("Failed adding to wishlist")
	//}
	return &dto.StatusMessageResponse{Status: !exists}, nil
}
