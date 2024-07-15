package story

import (
	"capstone/entities"
	"capstone/entities/doctor"
	"mime/multipart"
	"time"
)

type Story struct {
	Id       uint
	Title    string
	Content  string
	Date     time.Time
	ImageUrl string
	ViewCount int
	DoctorId uint
	Doctor   doctor.Doctor
	IsLiked bool
}

type RepositoryInterface interface {
	GetAllStories(metadata entities.Metadata, userId int, search string) ([]Story, error)
	GetStoryById(storyId int, userId int) (Story, error)
	GetLikedStories(metadata entities.Metadata, userId int) ([]Story, error)
	LikeStory(storyId int, userId int) error
	UnlikeStory(storyId int, userId int) error
	CountStoriesByDoctorId(doctorId int) (int, error)
	CountStoryLikesByDoctorId(doctorId int) (int, error)
	CountStoryViewByDoctorId(doctorId int) (int, error)
	CountStoryViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	PostStory(story Story) (Story, error)
	GetStoryByIdForDoctor(storyId int) (Story, error)
	GetAllStoriesByDoctorId(metadata entities.MetadataFull, doctorId int) ([]Story, error)
	EditStory(story Story) (Story, error)
	DeleteStory(storyId int) error
}

type UseCaseInterface interface {
	GetAllStories(metadata entities.Metadata, userId int, search string) ([]Story, error)
	GetStoryById(storyId int, userId int) (Story, error)
	GetLikedStories(metadata entities.Metadata, userId int) ([]Story, error)
	LikeStory(storyId int, userId int) error
	UnlikeStory(storyId int, userId int) error
	CountStoriesByDoctorId(doctorId int) (int, error)
	CountStoryLikesByDoctorId(doctorId int) (int, error)
	CountStoryViewByDoctorId(doctorId int) (int, error)
	CountStoryViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	PostStory(story Story, fileImage *multipart.FileHeader) (Story, error)
	GetStoryByIdForDoctor(storyId int) (Story, error)
	GetAllStoriesByDoctorId(metadata entities.MetadataFull, doctorId int) ([]Story, error)
	EditStory(story Story, fileImage *multipart.FileHeader) (Story, error)
	DeleteStory(storyId int) error
}