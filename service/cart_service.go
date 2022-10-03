package service

import (
	"kamari.com/config"
	"kamari.com/dto"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
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

	var response []dto.CartItemsResponse
	m := make(map[string][]dto.CartItemRaw)

	for _, item := range result {
		cartItem := dto.CartItemRaw{
			Id:              item.Id,
			Name:            item.Name,
			MainVariant:     item.MainVariant,
			MainVariantItem: item.MainVariantItem,
			SubVariant:      item.SubVariant,
			SubVariantItem:  item.SubVariantItem,
			PicUrl:          item.PicUrl,
			Price:           item.Price,
			Quantity:        item.Quantity,
			MerchantId:      item.MerchantId,
		}
		m[item.MerchantName] = append(m[item.MerchantName], cartItem)
	}

	for k, v := range m {
		var tempData []dto.CartItem
		for _, tempCartItem := range v {
			temp := dto.CartItem{
				Id:              tempCartItem.Id,
				Name:            tempCartItem.Name,
				MainVariant:     tempCartItem.MainVariant,
				MainVariantItem: tempCartItem.MainVariantItem,
				SubVariant:      tempCartItem.SubVariant,
				SubVariantItem:  tempCartItem.SubVariantItem,
				PicUrl:          tempCartItem.PicUrl,
				Price:           tempCartItem.Price,
				Quantity:        tempCartItem.Quantity,
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
