package forum

import (
	"capstone/constants"
	"capstone/entities"
	"capstone/entities/forum"
	userEntities "capstone/entities/user"
	"capstone/repositories/mysql/user"

	"gorm.io/gorm"
)

type ForumRepo struct {
	db *gorm.DB
}

func NewForumRepo(db *gorm.DB) *ForumRepo {
	return &ForumRepo{
		db: db,
	}
}

func (f *ForumRepo) JoinForum(forumId uint, userId uint) error {
	var forumMemberDB ForumMember
	forumMemberDB.ForumID = forumId
	forumMemberDB.UserID = userId

	err := f.db.Create(&forumMemberDB).Error
	if err != nil {
		return constants.ErrServer
	}
	return nil
}

func (f *ForumRepo) LeaveForum(forumId uint, userId uint) error {
	err := f.db.Where("forum_id = ? AND user_id = ?", forumId, userId).Delete(&ForumMember{}).Error
	if err != nil {
		return constants.ErrServer
	}
	return nil
}

func (f *ForumRepo) GetJoinedForum(userId uint, metadata entities.Metadata, search string) ([]forum.Forum, error) {
	var temps []ForumMember

	err := f.db.Where("user_id = ?", userId).Find(&temps).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	var forumIds []uint
	for _, temp := range temps {
		forumIds = append(forumIds, temp.ForumID)
	}

	var forumDBs []Forum

	query := f.db.Where("id IN ?", forumIds)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err = query.Find(&forumDBs).Error
	if err != nil {
		return nil, constants.ErrServer
	}
	
	counter := make([]int64, len(forumDBs))
	for i, forumDB := range forumDBs {
		err = f.db.Model(&ForumMember{}).Where("forum_id = ?", forumDB.ID).Count(&counter[i]).Error
		if err != nil {
			return nil, constants.ErrServer
		}
	}

	forumEnts := make([]forum.Forum, len(forumDBs))
	for i, forumDB := range forumDBs {
		forumEnts[i].ID = forumDB.ID
		forumEnts[i].Name = forumDB.Name
		forumEnts[i].ImageUrl = forumDB.ImageUrl
		forumEnts[i].NumberOfMembers = int(counter[i])
	}

	var forumMemberDBTemps []ForumMember
	for i:=0; i < len(forumEnts); i++ {
		err = f.db.Where("forum_id = ? AND user_id != ?", forumEnts[i].ID, userId).Limit(2).Find(&forumMemberDBTemps).Error
		if err != nil {
			return nil, constants.ErrServer
		}

		for _, forumMemberDBTemp := range forumMemberDBTemps {
			var userDB user.User
			err = f.db.Where("id = ?", forumMemberDBTemp.UserID).First(&userDB).Error
			if err != nil {
				return nil, constants.ErrServer
			}

			forumEnts[i].User = append(forumEnts[i].User, userEntities.User{
				Id:             userDB.Id,
				ProfilePicture: userDB.ProfilePicture,
			})
		}
	}

	return forumEnts, nil
}

func (f *ForumRepo) GetRecommendationForum(userId uint, metadata entities.Metadata, search string) ([]forum.Forum, error) {
	var forumMemberDB []ForumMember
	err := f.db.Preload("Forum").Where("user_id = ?", userId).Find(&forumMemberDB).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	var forumIDs []uint
	for _, forumMemberDBTemp := range forumMemberDB {
		forumIDs = append(forumIDs, forumMemberDBTemp.Forum.ID)
	}

	var query *gorm.DB
	var forumDB []Forum
	if len(forumIDs) != 0 {
		query = f.db.Model(&Forum{}).Where("id NOT IN ?", forumIDs).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit)
	} else {
		query = f.db.Model(&Forum{}).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit)
	}

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err = query.Find(&forumDB).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	counter := make([]int64, len(forumDB))
	for i, forumDBtemp := range forumDB {
		err = f.db.Model(&ForumMember{}).Where("forum_id = ?", forumDBtemp.ID).Count(&counter[i]).Error
		if err != nil {
			return nil, constants.ErrServer
		}
	}

	var forumEnts []forum.Forum
	for i, forumDBTemp := range forumDB {
		forumEnts = append(forumEnts, forum.Forum{
			ID:               forumDBTemp.ID,
			Name:             forumDBTemp.Name,
			ImageUrl:         forumDBTemp.ImageUrl,
			NumberOfMembers:  int(counter[i]),
		})
	}

	return forumEnts, nil
}

func (f *ForumRepo) GetForumById(forumId uint) (forum.Forum, error) {
	var forumDB Forum

	err := f.db.Where("id = ?", forumId).First(&forumDB).Error
	if err != nil {
		return forum.Forum{}, constants.ErrServer
	}

	var forumEnt forum.Forum
	forumEnt.ID = forumDB.ID
	forumEnt.Name = forumDB.Name
	forumEnt.Description = forumDB.Description
	forumEnt.ImageUrl = forumDB.ImageUrl

	return forumEnt, nil
}

func (f *ForumRepo) CreateForum(forumEnt forum.Forum) (forum.Forum, error) {
	var forumDB Forum
	forumDB.Name = forumEnt.Name
	forumDB.Description = forumEnt.Description
	forumDB.ImageUrl = forumEnt.ImageUrl
	forumDB.DoctorID = forumEnt.DoctorID

	err := f.db.Create(&forumDB).Error
	if err != nil {
		return forum.Forum{}, constants.ErrServer
	}

	return forum.Forum{
		ID:           forumDB.ID,
		Name:         forumDB.Name,
		Description:  forumDB.Description,
		ImageUrl:     forumDB.ImageUrl,
		DoctorID:     forumDB.DoctorID,
		NumberOfMembers: 0,
	}, nil
}

func (f *ForumRepo) GetAllForumsByDoctorId(doctorId uint, metadata entities.Metadata, search string) ([]forum.Forum, error) {
	var forumDBs []Forum
	query := f.db.Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Where("doctor_id = ?", doctorId)

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err := query.Find(&forumDBs).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	counter := make([]int64, len(forumDBs))
	for i, forumDB := range forumDBs {
		err = f.db.Model(&ForumMember{}).Where("forum_id = ?", forumDB.ID).Count(&counter[i]).Error
		if err != nil {
			return nil, constants.ErrServer
		}
	}

	var forumEnts []forum.Forum
	for i, forumDB := range forumDBs {
		forumEnts = append(forumEnts, forum.Forum{
			ID:               forumDB.ID,
			Name:             forumDB.Name,
			ImageUrl:         forumDB.ImageUrl,
			NumberOfMembers:  int(counter[i]),
		})
	}

	return forumEnts, nil
}

func (f *ForumRepo) UpdateForum(forumEnt forum.Forum) (forum.Forum, error) {
	var forumDB Forum

	err := f.db.Where("id = ?", forumEnt.ID).First(&forumDB).Error
	if err != nil {
		return forum.Forum{}, constants.ErrServer
	}

	forumDB.Name = forumEnt.Name
	forumDB.Description = forumEnt.Description

	if forumEnt.ImageUrl != "" {
		forumDB.ImageUrl = forumEnt.ImageUrl
	}

	err = f.db.Save(&forumDB).Error
	if err != nil {
		return forum.Forum{}, constants.ErrServer
	}

	return forum.Forum{
		ID:           forumDB.ID,
		Name:         forumDB.Name,
		Description:  forumDB.Description,
		ImageUrl:     forumDB.ImageUrl,
	}, nil
}

func (f *ForumRepo) DeleteForum(forumId uint) error {
	err := f.db.Where("id = ?", forumId).Delete(&Forum{}).Error
	if err != nil {
		return constants.ErrServer
	}
	return nil
}

func (f *ForumRepo) GetForumMemberByForumId(forumId uint, metadata entities.Metadata) ([]userEntities.User, error) {
	var forumMemberDBs []ForumMember
	err := f.db.Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Preload("User").Where("forum_id = ?", forumId).Find(&forumMemberDBs).Error
	if err != nil {
		return nil, constants.ErrServer
	}

	userEnts := make([]userEntities.User, len(forumMemberDBs))
	for i, forumMemberDB := range forumMemberDBs {
		userEnts[i] = userEntities.User{
			Id:       forumMemberDB.User.Id,
			Username: forumMemberDB.User.Username,
			Name:     forumMemberDB.User.Name,
			ProfilePicture: forumMemberDB.User.ProfilePicture,
		}
	}

	return userEnts, nil
}