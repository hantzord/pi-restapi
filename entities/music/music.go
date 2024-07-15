package music

import (
	"capstone/entities"
	"capstone/entities/doctor"
	"mime/multipart"
)

type Music struct {
	Id        uint
	Title     string
	Singer    string
	MusicUrl  string
	ImageUrl  string
	ViewCount int
	IsLiked   bool
	DoctorId  uint
	Doctor    doctor.Doctor
}

type RepositoryInterface interface {
	GetAllMusics(metadata entities.Metadata, userId int, search string) ([]Music, error)
	GetAllMusicsByDoctorId(metadata entities.MetadataFull, userId int) ([]Music, error)
	GetMusicById(musicId int, userId int) (Music, error)
	GetLikedMusics(metadata entities.Metadata, userId int) ([]Music, error)
	LikeMusic(musicId int, userId int) error
	UnlikeMusic(musicId int, userId int) error
	CountMusicByDoctorId(doctorId int) (int, error)
	CountMusicLikesByDoctorId(doctorId int) (int, error)
	CountMusicViewCountByDoctorId(doctorId int) (int, error)
	CountMusicViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	PostMusic(music Music) (Music, error)
	GetMusicByIdForDoctor(musicId int) (Music, error)
	EditMusic(music Music) (Music, error)
	DeleteMusic(musicId int) error
}

type UseCaseInterface interface {
	GetAllMusics(metadata entities.Metadata, userId int, search string) ([]Music, error)
	GetAllMusicsByDoctorId(metadata entities.MetadataFull, userId int) ([]Music, error)
	GetMusicById(musicId int, userId int) (Music, error)
	GetLikedMusics(metadata entities.Metadata, userId int) ([]Music, error)
	LikeMusic(musicId int, userId int) error
	UnlikeMusic(musicId int, userId int) error
	CountMusicByDoctorId(doctorId int) (int, error)
	CountMusicLikesByDoctorId(doctorId int) (int, error)
	CountMusicViewCountByDoctorId(doctorId int) (int, error)
	CountMusicViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	PostMusic(music Music, fileImage *multipart.FileHeader, fileMusic *multipart.FileHeader) (Music, error)
	GetMusicByIdForDoctor(musicId int) (Music, error)
	EditMusic(music Music, fileImage *multipart.FileHeader) (Music, error)
	DeleteMusic(musicId int) error
}