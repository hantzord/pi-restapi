package rating

import (
	"capstone/constants"
	"capstone/entities"
	ratingEntities "capstone/entities/rating"
	"capstone/entities/user"

	"gorm.io/gorm"
)

type RatingRepo struct {
	DB *gorm.DB
}

func NewRatingRepo(db *gorm.DB) *RatingRepo {
	return &RatingRepo{
		DB: db,
	}
}

func (repository *RatingRepo) SendFeedback(rating ratingEntities.Rating) (ratingEntities.Rating, error) {
	ratingDB := Rating{
		UserId:   rating.UserId,
		DoctorId: rating.DoctorId,
		Rate:     rating.Rate,
		Message:  rating.Message,
	}

	err := repository.DB.Create(&ratingDB).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}

	result := ratingEntities.Rating{
		Id:       ratingDB.ID,
		UserId:   ratingDB.UserId,
		DoctorId: ratingDB.DoctorId,
		Rate:     ratingDB.Rate,
		Message:  ratingDB.Message,
	}

	return result, nil
}

func (repository *RatingRepo) GetAllFeedbacks(metadata entities.Metadata, doctorId uint) ([]ratingEntities.Rating, error) {
	var ratingDB []Rating

	err := repository.DB.Limit(metadata.Limit).Offset((metadata.Page-1)*metadata.Limit).Preload("User").Where("doctor_id = ?", doctorId).Find(&ratingDB).Error
	if err != nil {
		return []ratingEntities.Rating{}, constants.ErrDataNotFound
	}

	result := make([]ratingEntities.Rating, len(ratingDB))
	for i := 0; i < len(ratingDB); i++ {
		result[i] = ratingEntities.Rating{
			Id:       ratingDB[i].ID,
			UserId:   ratingDB[i].UserId,
			User:     user.User{
				Id:       ratingDB[i].User.Id,
				Name:     ratingDB[i].User.Name,
				Username: ratingDB[i].User.Username,
				ProfilePicture: ratingDB[i].User.ProfilePicture,
			},
			Rate:     ratingDB[i].Rate,
			Message:  ratingDB[i].Message,
			Date:     ratingDB[i].CreatedAt.String(),
		}
	}

	return result, nil
}

func (repository *RatingRepo) GetSummaryRating(doctorId uint) (ratingEntities.Rating, error) {
	var oneStarCount int64
	var twoStarCount int64
	var threeStarCount int64
	var fourStarCount int64
	var fiveStarCount int64
	var average float64

	err := repository.DB.Model(&Rating{}).Where("doctor_id = ? AND rate = 5", doctorId).Count(&fiveStarCount).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}
	err = repository.DB.Model(&Rating{}).Where("doctor_id = ? AND rate = 4", doctorId).Count(&fourStarCount).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}
	err = repository.DB.Model(&Rating{}).Where("doctor_id = ? AND rate = 3", doctorId).Count(&threeStarCount).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}
	err = repository.DB.Model(&Rating{}).Where("doctor_id = ? AND rate = 2", doctorId).Count(&twoStarCount).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}
	err = repository.DB.Model(&Rating{}).Where("doctor_id = ? AND rate = 1", doctorId).Count(&oneStarCount).Error
	if err != nil {
		return ratingEntities.Rating{}, constants.ErrServer
	}

	total := (oneStarCount + twoStarCount + threeStarCount + fourStarCount + fiveStarCount)
	if total == 0 {
		average = 0
	}else{
		average = float64((oneStarCount + twoStarCount*2 + threeStarCount*3 + fourStarCount*4 + fiveStarCount*5) / total)
	}
	
	result := ratingEntities.Rating{
		OneStarCount:   int(oneStarCount),
		TwoStarCount:   int(twoStarCount),
		ThreeStarCount: int(threeStarCount),
		FourStarCount:  int(fourStarCount),
		FiveStarCount:  int(fiveStarCount),
		Average:        average,
	}

	return result, nil
}
