package service

import (
	"fmt"
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"time"
)

type CartService interface {
	GetCartHoverItems(user model.User) (*dto.CartHoverRes, error)
	GetCartItems(user model.User) ([]dto.CartItemsResponse, error)
	CreateCartItem(user model.User, req dto.AddCartItemRequest) error
	UpdateCartItem(id int, userId uint64, req dto.UpdateCartItemRequest) error
	DeleteCartItem(id int, userId uint64) error
}

type cartService struct {
	cartRepository repository.CartRepository
	appConfig      config.AppConfig
}

type CSConfig struct {
	CartRepository repository.CartRepository
	AppConfig      config.AppConfig
}

func NewCartService(c CSConfig) CartService {
	return cartService{
		cartRepository: c.CartRepository,
		appConfig:      c.AppConfig,
	}
}

func (c cartService) GetCartHoverItems(user model.User) (*dto.CartHoverRes, error) {
	result, err := c.cartRepository.GetCartHoverItems(user.Id)
	if err != nil {
		return nil, httperror.InternalServerError("Failed to get cart items")
	}
	return result, nil
}

func (c cartService) GetCartItems(user model.User) ([]dto.CartItemsResponse, error) {
	result, err := c.cartRepository.GetCartItems(user.Id)
	var discountedPrice uint64

	var response []dto.CartItemsResponse
	m := make(map[string][]dto.CartItemRaw)

	for _, item := range result {
		cartItem := dto.CartItemRaw{
			Id:              item.Id,
			VariantId:       item.VariantId,
			Name:            item.Name,
			MainVariant:     item.MainVariant,
			MainVariantItem: item.MainVariantItem,
			SubVariant:      item.SubVariant,
			SubVariantItem:  item.SubVariantItem,
			PicUrl:          item.PicUrl,
			Price:           item.Price,
			Quantity:        item.Quantity,
			MerchantId:      item.MerchantId,
			PromoTypeId:     item.PromoTypeId,
			DiscountAmount:  item.DiscountAmount,
			StartDate:       item.StartDate,
			EndDate:         item.EndDate,
			StockLeft:       item.StockLeft,
			MaximumItem:     item.MaximumItem,
			MinimumOrder:    item.MinimumOrder,
			MaximumOrder:    item.MaximumOrder,
		}
		m[item.MerchantName] = append(m[item.MerchantName], cartItem)
	}
	fmt.Println(m)

	for k, v := range m {
		var tempData []dto.CartItem
		for _, tempCartItem := range v {
			isValid := time.Now().Before(tempCartItem.EndDate) &&
				time.Now().After(tempCartItem.StartDate) &&
				tempCartItem.Quantity < tempCartItem.StockLeft &&
				tempCartItem.Quantity < tempCartItem.MaximumItem
			if !isValid {
				discountedPrice = tempCartItem.Price
			} else {
				if tempCartItem.PromoTypeId == 1 {
					discountedPrice = tempCartItem.Price - tempCartItem.DiscountAmount
				}
				if tempCartItem.PromoTypeId == 2 {
					discountedPrice = tempCartItem.DiscountAmount * (1 - tempCartItem.DiscountAmount)
				}
			}
			temp := dto.CartItem{
				Id:              tempCartItem.Id,
				VariantId:       tempCartItem.VariantId,
				Name:            tempCartItem.Name,
				MainVariant:     tempCartItem.MainVariant,
				MainVariantItem: tempCartItem.MainVariantItem,
				SubVariant:      tempCartItem.SubVariant,
				SubVariantItem:  tempCartItem.SubVariantItem,
				PicUrl:          tempCartItem.PicUrl,
				Price:           tempCartItem.Price,
				DiscountedPrice: discountedPrice,
				Quantity:        tempCartItem.Quantity,
				StockLeft:       tempCartItem.StockLeft,
				MinimumOrder:    tempCartItem.MinimumOrder,
				MaximumOrder:    tempCartItem.MaximumOrder,
			}
			tempData = append(tempData, temp)
		}
		response = append(response, dto.CartItemsResponse{
			Merchant: dto.CartItemMerchant{
				Id:   v[0].MerchantId,
				Name: k,
			},
			Items: tempData,
		})
	}

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c cartService) CreateCartItem(user model.User, req dto.AddCartItemRequest) error {
	err := c.cartRepository.CreateCartItem(user.Id, req)
	if err != nil {
		return err
	}
	return nil
}

func (c cartService) UpdateCartItem(id int, userId uint64, req dto.UpdateCartItemRequest) error {
	err := c.cartRepository.UpdateCartItem(id, userId, req)
	if err != nil {
		return err
	}
	return nil
}

func (c cartService) DeleteCartItem(id int, userId uint64) error {
	err := c.cartRepository.DeleteCartItem(id, userId)
	if err != nil {
		return err
	}
	return nil
}
