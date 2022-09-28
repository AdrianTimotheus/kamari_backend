package service

import (
	"kamari.com/dto"
	"kamari.com/repository"
)

type ProductReviewService interface {
	GetProductReviews(productId uint64) (*dto.ReviewResponse, error)
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

//TODO: add metada
func (p productReviewService) GetProductReviews(productId uint64) (*dto.ReviewResponse, error) {
	reviews, err := p.productReviewRepository.GetProductReview(productId)
	if err != nil {
		return nil, err
	}

	productReviews := make([]dto.ProductReviewResponse, len(reviews))
	
	reviewsResponse := &dto.ReviewResponse{
		ProductReviews: productReviews,
	}
	if len(reviews) < 1 {
		return reviewsResponse, nil
	}

	var score float32 = 0
	for i, review := range reviews {
		od, err := p.orderDetailRepository.GetOrderDetails(review.OrderDetailId)
		if err != nil {
			return nil, err
		}

		v, err := p.productVariantRepository.GetProductVariantById(od.ProductVariantId)
		if err != nil {
			return nil, err
		}

		r := dto.ProductReviewResponse{
			Id:             review.Id,
			Score:          review.Score,
			Feedback:       review.Feedback,
			ProductVariant: v,
			ReviewPhotos:   review.ReviewPhotos,
			CreatedAt:      review.CreatedAt,
		}
		productReviews[i] = r
		score += review.Score
	}

	reviewsResponse.ReviewCount = len(reviews)
	reviewsResponse.ReviewScore = score / float32(len(reviews))
	reviewsResponse.ProductReviews = productReviews

	return reviewsResponse, nil
}
