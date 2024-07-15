package forum

import (
	"capstone/entities"
	"capstone/entities/doctor"
	"capstone/entities/post"
	"capstone/entities/user"
	"mime/multipart"
)

type Forum struct {
	ID          	uint
	Name        	string
	Description 	string
	ImageUrl    	string
	DoctorID    	uint
	Doctor      	doctor.Doctor
	NumberOfMembers int
	User            []user.User
	Post            []post.Post
}

type ForumMember struct {
	ID      uint
	ForumID uint
	Forum   Forum
	UserID  uint
	User    user.User
}

type RepositoryInterface interface {
	JoinForum(forumId uint, userId uint) (error)
	LeaveForum(forumId uint, userId uint) (error)
	GetJoinedForum(userId uint, metadata entities.Metadata, search string) ([]Forum, error)
	GetRecommendationForum(userId uint, metadata entities.Metadata, search string) ([]Forum, error)
	GetForumById(forumId uint) (Forum, error)
	CreateForum(forum Forum) (Forum, error)
	GetAllForumsByDoctorId(doctorId uint, metadata entities.Metadata, search string) ([]Forum, error)
	UpdateForum(forum Forum) (Forum, error)
	DeleteForum(forumId uint) (error)
	GetForumMemberByForumId(forumId uint, metadata entities.Metadata) ([]user.User, error)
}

type UseCaseInterface interface {
	JoinForum(forumId uint, userId uint) (error)
	LeaveForum(forumId uint, userId uint) (error)
	GetJoinedForum(userId uint, metadata entities.Metadata, search string) ([]Forum, error)
	GetRecommendationForum(userId uint, metadata entities.Metadata, search string) ([]Forum, error)
	GetForumById(forumId uint) (Forum, error)
	CreateForum(forum Forum, fileImage *multipart.FileHeader) (Forum, error)
	GetAllForumsByDoctorId(doctorId uint, metadata entities.Metadata, search string) ([]Forum, error)
	UpdateForum(forum Forum, fileImage *multipart.FileHeader) (Forum, error)
	DeleteForum(forumId uint) (error)
	GetForumMemberByForumId(forumId uint, metadata entities.Metadata) ([]user.User, error)
}