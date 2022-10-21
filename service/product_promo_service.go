package service

import (
	"errors"
	"gorm.io/gorm"
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/httperror"
	"kamari.com/model"
	"kamari.com/repository"
	"time"
)

type ProductPromoService interface {
	CreateProductPromo(payload *dto.ProductPromoRequest, username string) error
	InsertPromoItems(payload []dto.PromoItems, promoId uint64) error
	GetProductPromos(merchantId uint64, keyword string, filters *helper.Filters) ([]model.ProductPromo, *helper.Metadata, error)
	GetProductPromoDetail(promoId uint64) (*dto.ProductPromoResponse, error)
	UpdateProductPromo(payload *dto.ProductPromoUpdate, username string, promoId uint64) error
	DeletePromoItem(merchantId, promoId uint64, itemId []uint64) error
	GetAvailableProductForPromo(merchantId uint64, keyword string, filters *helper.Filters) ([]dto.ProductResponse, *helper.Metadata, error)
}

type productPromoService struct {
	productPromoRepository repository.ProductPromoRepository
	productRepository      repository.ProductRepository
	merchantRepository     repository.MerchantRepository
}

type PPSConfig struct {
	ProductPromoRepository repository.ProductPromoRepository
	MerchantRepository     repository.MerchantRepository
	ProductRepository      repository.ProductRepository
}

func NewProductPromoService(c PPSConfig) ProductPromoService {
	return productPromoService{
		productPromoRepository: c.ProductPromoRepository,
		merchantRepository:     c.MerchantRepository,
		productRepository:      c.ProductRepository,
	}
}

func (p productPromoService) CreateProductPromo(payload *dto.ProductPromoRequest, username string) error {
	var (
		startDate time.Time
		endDate   time.Time
	)

	if payload.StartDate.Int64 > payload.EndDate.Int64 {
		return httperror.BadRequestError("end date cannot be smaller than start date", "INVALID_DATE")
	} else {
		startDate = time.Unix(payload.StartDate.Int64, 0)
		endDate = time.Unix(payload.EndDate.Int64, 0)
	}

	merchant, err := p.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		return err
	}

	productPromo := model.ProductPromo{
		Title:      payload.Title,
		StartDate:  startDate,
		EndDate:    endDate,
		MerchantId: merchant.Id,
	}

	promoItems := make([]model.ProductPromoItem, len(payload.PromoItems))
	for i, item := range payload.PromoItems {
		valid, err := p.productRepository.IsMerchantProduct(merchant.Id, item.ProductId)
		if err != nil || !valid {
			return errors.New("no data")
		}

		promoItems[i] = model.ProductPromoItem{
			DiscountAmount: item.DiscountAmount,
			MaximumItem:    item.MaximumItem,
			StockLeft:      item.StockLeft,
			ProductId:      item.ProductId,
			ProductPromoId: item.ProductPromoId,
			PromoTypeId:    item.PromoTypeId,
		}
	}

	if err = p.productPromoRepository.CreateMerchantPromoTransaction(&productPromo, promoItems); err != nil {
		return err
	}

	return nil
}

func (p productPromoService) InsertPromoItems(payload []dto.PromoItems, promoId uint64) error {
	promoItems := make([]model.ProductPromoItem, len(payload))
	for i, item := range payload {
		promoItems[i] = model.ProductPromoItem{
			ProductId:      item.ProductId,
			ProductPromoId: promoId,
			PromoTypeId:    item.PromoTypeId,
		}
	}

	err := p.productPromoRepository.InsertPromoItems(promoItems)
	if err != nil {
		return err
	}

	return nil
}

func (p productPromoService) GetProductPromos(merchantId uint64, keyword string, filters *helper.Filters) ([]model.ProductPromo, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 10)
	filters.Page = helper.ReadInt(filters.Page, 1)

	promos, metadata, err := p.productPromoRepository.GetProductPromos(merchantId, keyword, filters)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil, httperror.BadRequestError(err.Error(), "NO_DATA")
		default:
			return nil, nil, err
		}
	}

	return promos, metadata, nil
}

func (p productPromoService) GetProductPromoDetail(promoId uint64) (*dto.ProductPromoResponse, error) {
	promo, err := p.productPromoRepository.GetProductPromoDetail(promoId)
	if err != nil {
		return nil, err
	}

	pp := dto.ProductPromoResponse{
		Id:         promo.Id,
		Title:      promo.Title,
		StartDate:  promo.StartDate,
		EndDate:    promo.EndDate,
		MerchantId: promo.MerchantId,
	}

	promoItemResponse := make([]dto.ProductPromoItemsResponse, len(promo.ProductPromoItems))
	for i, item := range promo.ProductPromoItems {
		lowestPrice := helper.GetMinProductVarPrice(item.Product.ProductVariants)
		promoItemResponse[i] = dto.ProductPromoItemsResponse{
			Id:             item.Id,
			DiscountAmount: item.DiscountAmount,
			MaximumItem:    item.MaximumItem,
			StockLeft:      item.StockLeft,
			Price:          lowestPrice,
			ProductId:      item.ProductId,
			Product:        item.Product,
			ProductPromoId: item.ProductPromoId,
			PromoTypeId:    item.PromoTypeId,
			PromoType:      item.PromoType,
		}
	}

	pp.ProductPromoItems = promoItemResponse

	return &pp, nil
}

func (p productPromoService) DeletePromoItem(merchantId, promoId uint64, itemId []uint64) error {
	for _, id := range itemId {
		valid, err := p.productPromoRepository.IsValidPromoItem(merchantId, promoId, id)
		if err != nil || !valid {
			return httperror.ForbiddenError()
		}

		if removable, err := p.productPromoRepository.IsRemovable(promoId); err != nil || !removable {
			return err
		}
	}

	for _, id := range itemId {
		if err := p.productPromoRepository.DeletePromoItem(id); err != nil {
			return err
		}
	}

	return nil
}

func (p productPromoService) UpdateProductPromo(payload *dto.ProductPromoUpdate, username string, promoId uint64) error {
	var endDate time.Time

	if payload.EndDate.Valid {
		endDate = time.Unix(payload.EndDate.Int64, 0)
	}

	merchant, err := p.merchantRepository.GetMerchantDetail(username)
	if err != nil {
		return err
	}

	productPromo := model.ProductPromo{
		Id:      promoId,
		Title:   payload.Title,
		EndDate: endDate,
	}

	promoItems := make([]model.ProductPromoItem, len(payload.PromoItems))
	for i, item := range payload.PromoItems {
		valid, err := p.productRepository.IsMerchantProduct(merchant.Id, item.ProductId)
		if err != nil || !valid {
			return err
		}

		promoItems[i] = model.ProductPromoItem{
			DiscountAmount: item.DiscountAmount,
			MaximumItem:    item.MaximumItem,
			PromoTypeId:    item.PromoTypeId,
		}
	}

	if err = p.productPromoRepository.UpdatePromoTransaction(&productPromo, promoItems); err != nil {
		return err
	}

	return nil
}

func (p productPromoService) GetAvailableProductForPromo(merchantId uint64, keyword string, filters *helper.Filters) ([]dto.ProductResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 10)
	filters.Page = helper.ReadInt(filters.Page, 1)

	products, metadata, err := p.productPromoRepository.GetAvailableMerchantProductPromo(merchantId, keyword, filters)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		minPrice := helper.GetMinProductVarPrice(product.ProductVariants)
		res[i] = dto.ProductResponse{
			Id:           product.Id,
			Title:        product.Name,
			Price:        minPrice,
			ThumbnailUrl: product.ThumbnailUrl,
			Availability: product.Availability,
		}
	}

	return res, metadata, nil
}
