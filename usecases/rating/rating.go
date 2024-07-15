package rating

import (
	"capstone/constants"
	"capstone/entities"
	ratingEntities "capstone/entities/rating"
)

type RatingUseCase struct {
	ratingRepository ratingEntities.RepositoryInterface
}

func NewRatingUseCase(ratingRepository ratingEntities.RepositoryInterface) *RatingUseCase {
	return &RatingUseCase{
		ratingRepository: ratingRepository,
	}
}

func (ratingUseCase *RatingUseCase) SendFeedback(rating ratingEntities.Rating) (ratingEntities.Rating, error) {
	if rating.Rate < 1 || rating.Rate > 5 {
		return ratingEntities.Rating{}, constants.ErrInvalidRate
	}
	
	result, err := ratingUseCase.ratingRepository.SendFeedback(rating)
	if err != nil {
		return ratingEntities.Rating{}, err
	}
	return result, nil
}

func (ratingUseCase *RatingUseCase) GetAllFeedbacks(metadata entities.Metadata, doctorId uint) ([]ratingEntities.Rating, error) {
	result, err := ratingUseCase.ratingRepository.GetAllFeedbacks(metadata, doctorId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ratingUseCase *RatingUseCase) GetSummaryRating(doctorId uint) (ratingEntities.Rating, error) {
	result, err := ratingUseCase.ratingRepository.GetSummaryRating(doctorId)
	if err != nil {
		return ratingEntities.Rating{}, err
	}
	return result, nil
}