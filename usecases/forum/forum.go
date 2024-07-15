package forum

import (
	"capstone/constants"
	"capstone/entities"
	forumEntities "capstone/entities/forum"
	userEntities "capstone/entities/user"
	"capstone/utilities"
	"mime/multipart"
)

type ForumUseCase struct {
	forumInterface forumEntities.RepositoryInterface
}

func NewForumUseCase(forumInterface forumEntities.RepositoryInterface) *ForumUseCase {
	return &ForumUseCase{
		forumInterface: forumInterface,
	}
}

func (forumUseCase *ForumUseCase) JoinForum(forumId uint, userId uint) error {
	if forumId == 0 {
		return constants.ErrEmptyInputForum
	}

	err := forumUseCase.forumInterface.JoinForum(forumId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (forumUseCase *ForumUseCase) GetJoinedForum(userId uint, metadata entities.Metadata, search string) ([]forumEntities.Forum, error) {
	forums, err := forumUseCase.forumInterface.GetJoinedForum(userId, metadata, search)
	if err != nil {
		return nil, err
	}
	return forums, nil
}

func (forumUseCase *ForumUseCase) LeaveForum(forumId uint, userId uint) error {
	if forumId == 0 {
		return constants.ErrEmptyInputForum
	}

	err := forumUseCase.forumInterface.LeaveForum(forumId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (forumUseCase *ForumUseCase) GetRecommendationForum(userId uint, metadata entities.Metadata, search string) ([]forumEntities.Forum, error) {
	forums, err := forumUseCase.forumInterface.GetRecommendationForum(userId, metadata, search)
	if err != nil {
		return nil, err
	}
	return forums, nil
}

func (forumUseCase *ForumUseCase) GetForumById(forumId uint) (forumEntities.Forum, error) {
	forum, err := forumUseCase.forumInterface.GetForumById(forumId)
	if err != nil {
		return forumEntities.Forum{}, err
	}
	return forum, nil
}

func (forumUseCase *ForumUseCase) CreateForum(forum forumEntities.Forum, fileImage *multipart.FileHeader) (forumEntities.Forum, error) {
	if forum.Name == "" || forum.Description == "" {
		return forumEntities.Forum{}, constants.ErrEmptyInputForum
	}

	if fileImage != nil {
		secureUrl, err := utilities.UploadImage(fileImage)
		if err != nil {
			return forumEntities.Forum{}, constants.ErrUploadImage
		}
		forum.ImageUrl = secureUrl
	}

	forum, err := forumUseCase.forumInterface.CreateForum(forum)
	if err != nil {
		return forumEntities.Forum{}, err
	}
	return forum, nil
}

func (forumUseCase *ForumUseCase) GetAllForumsByDoctorId(doctorId uint, metadata entities.Metadata, search string) ([]forumEntities.Forum, error) {
	forums, err := forumUseCase.forumInterface.GetAllForumsByDoctorId(doctorId, metadata, search)
	if err != nil {
		return nil, err
	}
	return forums, nil
}

func (forumUseCase *ForumUseCase) UpdateForum(forum forumEntities.Forum, fileImage *multipart.FileHeader) (forumEntities.Forum, error) {
	if forum.Name == "" || forum.Description == "" {
		return forumEntities.Forum{}, constants.ErrEmptyCreateForum
	}

	if fileImage != nil {
		secureUrl, err := utilities.UploadImage(fileImage)
		if err != nil {
			return forumEntities.Forum{}, constants.ErrUploadImage
		}
		forum.ImageUrl = secureUrl
	}else{
		forum.ImageUrl = ""
	}

	forum, err := forumUseCase.forumInterface.UpdateForum(forum)
	if err != nil {
		return forumEntities.Forum{}, err
	}
	return forum, nil
}

func (forumUseCase *ForumUseCase) DeleteForum(forumId uint) error {
	err := forumUseCase.forumInterface.DeleteForum(forumId)
	if err != nil {
		return err
	}
	return nil
}

func (forumUseCase *ForumUseCase) GetForumMemberByForumId(forumId uint, metadata entities.Metadata) ([]userEntities.User, error) {
	members, err := forumUseCase.forumInterface.GetForumMemberByForumId(forumId, metadata)
	if err != nil {
		return nil, err
	}
	return members, nil
}