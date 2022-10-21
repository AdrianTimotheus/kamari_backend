package service

import (
	"kamari.com/dto"
	"kamari.com/helper"
	"kamari.com/repository"
)

type ProductReviewService interface {
	GetProductReviews(productId uint64, filters *helper.Filters) ([]dto.ProductReviewResponse, *helper.Metadata, error)
	CreateReview(orderDetailId uint64, userId uint64, req dto.CreateReviewRequest) error
}

type productReviewService struct {
	productReviewRepository  repository.ProductReviewRepository
	orderDetailRepository    repository.OrderDetailsRepository
	productVariantRepository repository.ProductVariantRepository
}

type PRSConfig struct {
	ProductReviewRepository  repository.ProductReviewRepository
	OrderDetailRepository    repository.OrderDetailsRepository
	ProductVariantRepository repository.ProductVariantRepository
}

func NewProductReviewService(c PRSConfig) ProductReviewService {
	return productReviewService{
		productReviewRepository:  c.ProductReviewRepository,
		orderDetailRepository:    c.OrderDetailRepository,
		productVariantRepository: c.ProductVariantRepository,
	}
}

func (p productReviewService) GetProductReviews(productId uint64, filters *helper.Filters) ([]dto.ProductReviewResponse, *helper.Metadata, error) {
	filters.Limit = helper.ReadInt(filters.Limit, 6)
	filters.Page = helper.ReadInt(filters.Page, 1)

	reviews, metadata, err := p.productReviewRepository.GetProductReview(productId, filters)
	if err != nil {
		return nil, nil, err
	}

	if filters.FilterBy == "withMedia" {
		reviews, metadata = helper.FilterReviewsByMedia(reviews, metadata)
	}

	productReviews := make([]dto.ProductReviewResponse, len(reviews))

	if len(reviews) < 1 {
		return productReviews, nil, nil
	}

	var score float32 = 0
	for i, review := range reviews {
		od, err := p.orderDetailRepository.GetOrderDetails(review.OrderDetailId)
		if err != nil {
			return nil, nil, err
		}

		rev, err := p.orderDetailRepository.GetOrdererUser(review.OrderDetailId)
		if err != nil {
			return nil, nil, err
		}

		reviewer := dto.NewReviewer(rev.Username, rev.PicUrl)

		v, err := p.productVariantRepository.GetProductVariantById(od.ProductVariantId)
		if err != nil {
			return nil, nil, err
		}

		r := dto.ProductReviewResponse{
			Id:             review.Id,
			Reviewer:       reviewer,
			Score:          review.Score,
			Feedback:       review.Feedback,
			ProductVariant: v,
			ReviewPhotos:   review.ReviewPhotos,
			CreatedAt:      review.CreatedAt,
		}
		productReviews[i] = r
		score += review.Score
	}

	return productReviews, metadata, nil
}

func (p productReviewService) CreateReview(orderDetailId uint64, userId uint64, req dto.CreateReviewRequest) error {
	err := p.productReviewRepository.CreateReview(orderDetailId, userId, req)
	if err != nil {
		return err
	}
	return nil
}
