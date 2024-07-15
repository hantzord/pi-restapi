package stories

import (
	"capstone/constants"
	"capstone/entities"
	storyEntities "capstone/entities/story"
	"capstone/utilities"
	"mime/multipart"
)

type StoryUseCase struct {
	storyRepository storyEntities.RepositoryInterface
}

func NewStoryUseCase(storyRepository storyEntities.RepositoryInterface) *StoryUseCase {
	return &StoryUseCase{
		storyRepository: storyRepository,
	}
}

func (storiesUseCase *StoryUseCase) GetAllStories(metadata entities.Metadata, userId int, search string) ([]storyEntities.Story, error) {
	stories, err := storiesUseCase.storyRepository.GetAllStories(metadata, userId, search)
	if err != nil {
		return []storyEntities.Story{}, err
	}
	return stories, nil
}

func (storiesUseCase *StoryUseCase) GetStoryById(storyId int, userId int) (storyEntities.Story, error) {
	story, err := storiesUseCase.storyRepository.GetStoryById(storyId, userId)
	if err != nil {
		return storyEntities.Story{}, err
	}
	return story, nil
}

func (storiesUseCase *StoryUseCase) GetLikedStories(metadata entities.Metadata, userId int) ([]storyEntities.Story, error) {
	stories, err := storiesUseCase.storyRepository.GetLikedStories(metadata, userId)
	if err != nil {
		return []storyEntities.Story{}, err
	}
	return stories, nil
}

func (storiesUseCase *StoryUseCase) LikeStory(storyId int, userId int) error {
	err := storiesUseCase.storyRepository.LikeStory(storyId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (storiesUseCase *StoryUseCase) UnlikeStory(storyId int, userId int) error {
	err := storiesUseCase.storyRepository.UnlikeStory(storyId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (storiesUseCase *StoryUseCase) CountStoriesByDoctorId(doctorId int) (int, error) {
	count, err := storiesUseCase.storyRepository.CountStoriesByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (storiesUseCase *StoryUseCase) CountStoryLikesByDoctorId(doctorId int) (int, error) {
	count, err := storiesUseCase.storyRepository.CountStoryLikesByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (storiesUseCase *StoryUseCase) CountStoryViewByDoctorId(doctorId int) (int, error) {
	count, err := storiesUseCase.storyRepository.CountStoryViewByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (storiesUseCase *StoryUseCase) CountStoryViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error) {
	if startMonth == "" || endMonth == "" {
		return map[int]int{}, constants.ErrEmptyInputViewByMonth
	}
	
	count, err := storiesUseCase.storyRepository.CountStoryViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return map[int]int{}, err
	}
	return count, nil
}

func (storiesUseCase *StoryUseCase) PostStory(story storyEntities.Story, file *multipart.FileHeader) (storyEntities.Story, error) {
	if story.Title == "" || story.Content == "" {
		return storyEntities.Story{}, constants.ErrEmptyInputStory
	}

	if file != nil {
		secureUrl, err := utilities.UploadImage(file)
		if err != nil {
			return storyEntities.Story{}, constants.ErrUploadImage
		}
		story.ImageUrl = secureUrl
	}

	story, err := storiesUseCase.storyRepository.PostStory(story)
	if err != nil {
		return storyEntities.Story{}, err
	}
	return story, nil
}

func (storiesUseCase *StoryUseCase) GetStoryByIdForDoctor(storyId int) (storyEntities.Story, error) {
	story, err := storiesUseCase.storyRepository.GetStoryByIdForDoctor(storyId)
	if err != nil {
		return storyEntities.Story{}, err
	}
	return story, nil
}

func (storiesUseCase *StoryUseCase) GetAllStoriesByDoctorId(metadata entities.MetadataFull, doctorId int) ([]storyEntities.Story, error) {
	stories, err := storiesUseCase.storyRepository.GetAllStoriesByDoctorId(metadata, doctorId)
	if err != nil {
		return []storyEntities.Story{}, err
	}
	return stories, nil
}

func (storiesUseCase *StoryUseCase) EditStory(story storyEntities.Story, file *multipart.FileHeader) (storyEntities.Story, error) {
	if story.Title == "" || story.Content == "" {
		return storyEntities.Story{}, constants.ErrEmptyInputStory
	}
	if file != nil {
		secureUrl, err := utilities.UploadImage(file)
		if err != nil {
			return storyEntities.Story{}, constants.ErrUploadImage
		}
		story.ImageUrl = secureUrl
	}else{
		story.ImageUrl = ""
	}

	story, err := storiesUseCase.storyRepository.EditStory(story)
	if err != nil {
		return storyEntities.Story{}, err
	}
	return story, nil
}

func (storiesUseCase *StoryUseCase) DeleteStory(storyId int) error {
	err := storiesUseCase.storyRepository.DeleteStory(storyId)
	if err != nil {
		return err
	}
	return nil
}