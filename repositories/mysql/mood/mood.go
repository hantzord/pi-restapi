package mood

import (
	"capstone/constants"
	moodEntities "capstone/entities/mood"

	"gorm.io/gorm"
)

type MoodRepo struct {
	db *gorm.DB
}

func NewMoodRepo(db *gorm.DB) *MoodRepo {
	return &MoodRepo{
		db: db,
	}
}

func (moodRepo *MoodRepo) SendMood(mood moodEntities.Mood) (moodEntities.Mood, error) {
	moodDB := Mood{
		UserId:     mood.UserId,
		MoodTypeId: mood.MoodTypeId,
		Message:    mood.Message,
		Date:       mood.Date,
		ImageUrl:   mood.ImageUrl,
	}

	err := moodRepo.db.Create(&moodDB).Error
	if err != nil {
		return moodEntities.Mood{}, constants.ErrServer
	}

	err = moodRepo.db.Model(&moodDB).Preload("MoodType").Find(&moodDB).Error
	if err != nil {
		return moodEntities.Mood{}, constants.ErrServer
	}

	result := moodEntities.Mood{
		ID:         moodDB.ID,
		UserId:     moodDB.UserId,
		MoodTypeId: moodDB.MoodTypeId,
		MoodType:   moodEntities.MoodType{
			ID:   moodDB.MoodType.ID,
			Name: moodDB.MoodType.Name,
		},
		Message:    moodDB.Message,
		Date:       moodDB.Date,
		ImageUrl:   moodDB.ImageUrl,
	}

	return result, nil
}

func (moodRepo *MoodRepo) GetAllMoods(userId int, startDate string, endDate string) ([]moodEntities.Mood, error) {
	var moodsDB []Mood
	err := moodRepo.db.Preload("MoodType").Where("user_id = ? AND date BETWEEN ? AND ?", userId, startDate, endDate).Find(&moodsDB).Error
	if err != nil {
		return []moodEntities.Mood{}, constants.ErrServer
	}

	moodEnts := make([]moodEntities.Mood, len(moodsDB))

	for i, moodDB := range moodsDB {
		date := moodDB.Date
		if len(date) > 10 {
            date = date[:10]
        }

		moodEnts[i] = moodEntities.Mood{
			ID:         moodDB.ID,
			MoodType:   moodEntities.MoodType{
				ID:   moodDB.MoodType.ID,
				Name: moodDB.MoodType.Name,
			},
			Date:       date,
		}
	}

	return moodEnts, nil
}

func (moodRepo *MoodRepo) GetMoodById(moodId int) (moodEntities.Mood, error) {
	var moodDB Mood
	err := moodRepo.db.Preload("MoodType").Where("id = ?", moodId).Find(&moodDB).Error
	if err != nil {
		return moodEntities.Mood{}, constants.ErrServer
	}

	result := moodEntities.Mood{
		ID:         moodDB.ID,
		MoodType:   moodEntities.MoodType{
			ID:   moodDB.MoodType.ID,
			Name: moodDB.MoodType.Name,
		},
		Message:    moodDB.Message,
		Date:       moodDB.Date,
		ImageUrl:   moodDB.ImageUrl,
	}

	return result, nil
}