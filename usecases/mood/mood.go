package mood

import (
	"capstone/constants"
	moodEntities "capstone/entities/mood"
	"capstone/utilities"
	"mime/multipart"
	"time"
)

type MoodUseCase struct {
	moodInterface moodEntities.RepositoryInterface
}

func NewMoodUseCase(moodInterface moodEntities.RepositoryInterface) *MoodUseCase {
	return &MoodUseCase{
		moodInterface: moodInterface,
	}
}

func (moodUseCase *MoodUseCase) SendMood(file *multipart.FileHeader, mood moodEntities.Mood) (moodEntities.Mood, error) {
	if mood.MoodTypeId == 0 || mood.Date == "" {
		return moodEntities.Mood{}, constants.ErrEmptyInputMood
	}

	if file != nil {
		secureUrl, err := utilities.UploadImage(file)
		if err != nil {
			return moodEntities.Mood{}, constants.ErrUploadImage
		}
		mood.ImageUrl = secureUrl
	}

	mood, err := moodUseCase.moodInterface.SendMood(mood)
	if err != nil {
		return moodEntities.Mood{}, err
	}
	
	return mood, nil
}

func (moodUseCase *MoodUseCase) GetAllMoods(userId int, startDateStr string, endDateStr string) ([]moodEntities.Mood, error) {
	if startDateStr == "" || endDateStr == "" {
		return []moodEntities.Mood{}, constants.ErrEmptyRangeDateMood
	}
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return []moodEntities.Mood{}, constants.ErrInvalidStartDate
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return []moodEntities.Mood{}, constants.ErrInvalidEndDate
	}

	if startDate.After(endDate) {
        return []moodEntities.Mood{}, constants.ErrStartDateGreater
    }

	result, err := moodUseCase.moodInterface.GetAllMoods(userId, startDateStr, endDateStr)
	if err != nil {
		return []moodEntities.Mood{}, err
	}
	return result, nil
}

func (moodUseCase *MoodUseCase) GetMoodById(moodId int) (moodEntities.Mood, error) {
	result, err := moodUseCase.moodInterface.GetMoodById(moodId)
	if err != nil {
		return moodEntities.Mood{}, err
	}
	return result, nil
}