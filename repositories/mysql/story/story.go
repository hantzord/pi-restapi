package story

import (
	"capstone/constants"
	"capstone/entities"
	doctorEntities "capstone/entities/doctor"
	storyEntities "capstone/entities/story"

	"gorm.io/gorm"
)

type StoriesRepo struct {
	DB *gorm.DB
}

func NewStoryRepo(db *gorm.DB) *StoriesRepo {
	return &StoriesRepo{
		DB: db,
	}
}

func (repository *StoriesRepo) GetAllStories(metadata entities.Metadata, userId int, search string) ([]storyEntities.Story, error) {
	var storiesDb []Story

	query := repository.DB.Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Preload("Doctor")

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	err := query.Find(&storiesDb).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	storyLikes := make([]StoryLikes, len(storiesDb))
	var counter int64
	var isLiked []bool

	for i := 0; i < len(storiesDb); i++ {
		storyLikes[i].UserId = uint(userId)
		storyLikes[i].StoryId = storiesDb[i].ID
		err = repository.DB.Model(&storyLikes[i]).Where("user_id = ? AND story_id = ?", storyLikes[i].UserId, storyLikes[i].StoryId).Count(&counter).Error

		if err != nil {
			return nil, constants.ErrServer
		}

		if counter > 0 {
			isLiked = append(isLiked, true)
		} else {
			isLiked = append(isLiked, false)
		}

		counter = 0
	}

	storiesEnt := make([]storyEntities.Story, len(storiesDb))
	for i := 0; i < len(storiesDb); i++ {
		storiesEnt[i] = storyEntities.Story{
			Id:        storiesDb[i].ID,
			Title:     storiesDb[i].Title,
			Content:   storiesDb[i].Content,
			Date:      storiesDb[i].Date,
			ImageUrl:  storiesDb[i].ImageUrl,
			DoctorId:  storiesDb[i].DoctorId,
			Doctor: doctorEntities.Doctor{
				ID:   storiesDb[i].Doctor.ID,
				Name: storiesDb[i].Doctor.Name,
			},
			IsLiked: isLiked[i],
		}
	}

	return storiesEnt, nil
}

func (repository *StoriesRepo) GetStoryById(storyId int, userId int) (storyEntities.Story, error) {
	var storyDb Story
	err := repository.DB.Where("id = ?", storyId).Preload("Doctor").First(&storyDb).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrDataNotFound
	}

	var storyLikes StoryLikes
	var isLiked bool
	var counter int64

	err = repository.DB.Model(&storyLikes).Where("user_id = ? AND story_id = ?", userId, storyId).Count(&counter).Error

	if err != nil {
		return storyEntities.Story{}, constants.ErrServer
	}

	if counter > 0 {
		isLiked = true
	} else {
		isLiked = false
	}

	storyResp := storyEntities.Story{
		Id:        storyDb.ID,
		Title:     storyDb.Title,
		Content:   storyDb.Content,
		Date:      storyDb.Date,
		ImageUrl:  storyDb.ImageUrl,
		DoctorId:  storyDb.DoctorId,
		Doctor: doctorEntities.Doctor{
			ID:   storyDb.Doctor.ID,
			Name: storyDb.Doctor.Name,
		},
		IsLiked: isLiked,
	}

	err = repository.DB.Model(&StoryViews{}).Create(&StoryViews{UserId: uint(userId), StoryId: uint(storyId)}).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrServer
	}

	return storyResp, nil
}

func (repository *StoriesRepo) GetLikedStories(metadata entities.Metadata, userId int) ([]storyEntities.Story, error) {
	var storyLikesDb []StoryLikes
	err := repository.DB.Limit(metadata.Limit).Offset((metadata.Page-1)*metadata.Limit).Where("user_id = ?", userId).Find(&storyLikesDb).Error
	if err != nil {
		return nil, constants.ErrDataNotFound
	}

	var likedStoryIDs []int
	for i := 0; i < len(storyLikesDb); i++ {
		likedStoryIDs = append(likedStoryIDs, int(storyLikesDb[i].StoryId))
	}

	var storiesDb []Story
	err = repository.DB.Where("id IN ?", likedStoryIDs).Preload("Doctor").Find(&storiesDb).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	storiesEnt := make([]storyEntities.Story, len(storiesDb))
	for i := 0; i < len(storiesDb); i++ {
		storiesEnt[i] = storyEntities.Story{
			Id:        storiesDb[i].ID,
			Title:     storiesDb[i].Title,
			Content:   storiesDb[i].Content,
			Date:      storiesDb[i].Date,
			ImageUrl:  storiesDb[i].ImageUrl,
			DoctorId:  storiesDb[i].DoctorId,
			Doctor: doctorEntities.Doctor{
				ID:   storiesDb[i].Doctor.ID,
				Name: storiesDb[i].Doctor.Name,
			},
			IsLiked: true,
		}
	}

	return storiesEnt, nil
}

func (repository *StoriesRepo) LikeStory(storyId int, userId int) error {
	var storyLikesDb StoryLikes

	err := repository.DB.Where("user_id = ? AND story_id = ?", userId, storyId).First(&storyLikesDb).Error
	if err == nil {
		return constants.ErrAlreadyLiked
	}

	storyLikesDb.UserId = uint(userId)
	storyLikesDb.StoryId = uint(storyId)

	err = repository.DB.Create(&storyLikesDb).Error
	if err != nil {
		return constants.ErrServer
	}

	return nil
}

func (repository *StoriesRepo) UnlikeStory(storyId int, userId int) error {
	err := repository.DB.Where("user_id = ? AND story_id = ?", userId, storyId).Delete(&StoryLikes{}).Error
	if err != nil {
		return constants.ErrServer
	}

	return nil
}

func (repository *StoriesRepo) CountStoriesByDoctorId(doctorId int) (int, error) {
	var counter int64
	err := repository.DB.Model(&Story{}).Where("doctor_id = ?", doctorId).Count(&counter).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(counter), nil
}

func (repository *StoriesRepo) CountStoryLikesByDoctorId(doctorId int) (int, error) {
	var counter int64
	err := repository.DB.Table("story_likes").
		Joins("JOIN stories ON story_likes.story_id = stories.id").
		Where("stories.doctor_id = ?", doctorId).
		Count(&counter).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(counter), nil
}

func (repository *StoriesRepo) CountStoryViewByDoctorId(doctorId int) (int, error) {
	var storyDb []Story
	err := repository.DB.Model(&Story{}).Where("doctor_id = ?", doctorId).Find(&storyDb).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	var storyDBIDs []int
	for _, story := range storyDb {
		storyDBIDs = append(storyDBIDs, int(story.ID))
	}
	
	var totalViews int64
	err = repository.DB.Model(&StoryViews{}).Where("story_id IN ?", storyDBIDs).Count(&totalViews).Error
	if err != nil {
		return 0, constants.ErrServer
	}

	return int(totalViews), nil
}

func (repository *StoriesRepo) CountStoryViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error) {
	var storyDB []Story
	err := repository.DB.Model(&Story{}).Where("doctor_id = ?", doctorId).Find(&storyDB).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	var StoryDBIDs []int
	for _, story := range storyDB {
		StoryDBIDs = append(StoryDBIDs, int(story.ID))
	}

	if len(StoryDBIDs) == 0 {
		return nil, constants.ErrDataNotFound
	}

	var results []struct {
        Month int
        Views int
    }

	query := repository.DB.Model(&StoryViews{}).Select("MONTH(created_at) as month, COUNT(*) as views").
        Where("story_id IN ?", StoryDBIDs).
        Where("created_at BETWEEN ? AND ?", startMonth+"-01", endMonth+"-31").
        Where("deleted_at IS NULL").
        Group("month").
        Order("month")

    err = query.Scan(&results).Error
    if err != nil {
        return nil, constants.ErrServer
    }

    viewsByMonth := make(map[int]int)
    for _, result := range results {
        viewsByMonth[result.Month] = result.Views
    }

    return viewsByMonth, nil
}

func (repository *StoriesRepo) PostStory(story storyEntities.Story) (storyEntities.Story, error) {
	storyDb := Story{
		Title:     story.Title,
		Content:   story.Content,
		Date:      story.Date,
		ImageUrl:  story.ImageUrl,
		DoctorId:  story.DoctorId,
	}

	err := repository.DB.Create(&storyDb).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrServer
	}

	return storyEntities.Story{
		Id:        storyDb.ID,
		Title:     storyDb.Title,
		Content:   storyDb.Content,
		Date:      storyDb.Date,
		ImageUrl:  storyDb.ImageUrl,
		DoctorId:  storyDb.DoctorId,
	}, nil
}

func (repository *StoriesRepo) GetStoryByIdForDoctor(storyId int) (storyEntities.Story, error) {
	var storyDb Story
	err := repository.DB.Where("id = ?", storyId).First(&storyDb).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrDataNotFound
	}

	var totalView int64
	err = repository.DB.Model(&StoryViews{}).Where("story_id = ?", storyId).Count(&totalView).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrServer
	}

	return storyEntities.Story{
		Id:        storyDb.ID,
		Title:     storyDb.Title,
		Content:   storyDb.Content,
		Date:      storyDb.Date,
		ImageUrl:  storyDb.ImageUrl,
		DoctorId:  storyDb.DoctorId,
		ViewCount: int(totalView),
	}, nil
}

func (repository *StoriesRepo) GetAllStoriesByDoctorId(metadata entities.MetadataFull, doctorId int) ([]storyEntities.Story, error) {
	var storiesDb []Story

	query := repository.DB.Where("doctor_id = ?", doctorId).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Order(metadata.Sort + " " + metadata.Order)

	if metadata.Search != "" {
		query = query.Where("title LIKE ?", "%"+metadata.Search+"%")
	}

	err := query.Find(&storiesDb).Error
	if err != nil {
		return []storyEntities.Story{}, constants.ErrServer
	}

	var storyDBIDs []int
	for i := 0; i < len(storiesDb); i++ {
		storyDBIDs = append(storyDBIDs, int(storiesDb[i].ID))
	}

	var totalViews []int
	for i := 0; i < len(storyDBIDs); i++ {
		var totalView int64
		err = repository.DB.Model(&StoryViews{}).Where("story_id = ?", storyDBIDs[i]).Count(&totalView).Error
		if err != nil {
			return []storyEntities.Story{}, constants.ErrServer
		}
		totalViews = append(totalViews, int(totalView))
	}

	storiesEnt := make([]storyEntities.Story, len(storiesDb))
	for i := 0; i < len(storiesDb); i++ {
		storiesEnt[i] = storyEntities.Story{
			Id:        storiesDb[i].ID,
			Title:     storiesDb[i].Title,
			Content:   storiesDb[i].Content,
			Date:      storiesDb[i].Date,
			ImageUrl:  storiesDb[i].ImageUrl,
			DoctorId:  storiesDb[i].DoctorId,
			ViewCount: totalViews[i],
		}
	}

	return storiesEnt, nil
}

func (repository *StoriesRepo) EditStory(story storyEntities.Story) (storyEntities.Story, error) {
	var storyDB Story

	err := repository.DB.Where("id = ?", story.Id).First(&storyDB).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrDataNotFound
	}

	storyDB.Title = story.Title
	storyDB.Content = story.Content

	if story.ImageUrl != "" {
		storyDB.ImageUrl = story.ImageUrl
	}

	err = repository.DB.Save(&storyDB).Error
	if err != nil {
		return storyEntities.Story{}, constants.ErrServer
	}

	return storyEntities.Story{
		Id:        storyDB.ID,
		Title:     storyDB.Title,
		Content:   storyDB.Content,
		Date:      storyDB.Date,
		ImageUrl:  storyDB.ImageUrl,
		DoctorId:  storyDB.DoctorId,
	}, nil
}

func (repository *StoriesRepo) DeleteStory(storyId int) error {
	return repository.DB.Where("id = ?", storyId).Delete(&Story{}).Error
}