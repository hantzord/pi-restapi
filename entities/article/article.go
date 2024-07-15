package article

import (
	"capstone/controllers/article/response"
	"capstone/entities"
	"capstone/entities/doctor"
	"mime/multipart"
	"time"
)

type Article struct {
	ID          uint
	Title       string
	Content     string
	Date        time.Time
	ImageUrl    string
	ViewCount   int
	DoctorID    uint
	Doctor      doctor.Doctor
	IsLiked     bool
	ReadingTime int
}

type ArticleRepositoryInterface interface {
	CreateArticle(article *Article, userId int) (*Article, error)
	GetAllArticle(metadata entities.Metadata, userId int, search string) ([]Article, error)
	GetArticleById(articleId int, userId int) (Article, error)
	GetLikedArticle(metadata entities.Metadata, userId int) ([]Article, error)
	LikeArticle(articleId int, userId int) error
	UnlikeArticle(articleId int, userId int) error
	GetArticleByIdForDoctor(articleId int) (Article, error)
	GetAllArticleByDoctorId(metadata entities.MetadataFull, doctorId int) ([]Article, error)
	CountArticleByDoctorId(doctorId int) (int, error)
	CountArticleLikesByDoctorId(doctorId int) (int, error)
	CountArticleViewByDoctorId(doctorId int) (int, error)
	CountArticleViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	EditArticle(article Article) (Article, error)
	DeleteArticle(articleId int) error
	// IncrementViewCount(articleId int) error
}

type ArticleUseCaseInterface interface {
	CreateArticle(article *Article, userId int) (*Article, error)
	GetAllArticle(metadata entities.Metadata, userId int, search string) ([]Article, error)
	GetArticleById(articleId int, userId int) (Article, error)
	GetLikedArticle(metadata entities.Metadata, userId int) ([]Article, error)
	LikeArticle(articleId int, userId int) error
	UnlikeArticle(articleId int, userId int) error
	GetArticleByIdForDoctor(articleId int) (Article, error)
	GetAllArticleByDoctorId(metadata entities.MetadataFull, doctorId int) ([]Article, error)
	CountArticleByDoctorId(doctorId int) (int, error)
	CountArticleLikesByDoctorId(doctorId int) (int, error)
	CountArticleViewByDoctorId(doctorId int) (int, error)
	CountArticleViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error)
	EditArticle(article Article, file *multipart.FileHeader) (Article, error)
	DeleteArticle(articleId int) error
}

func (ar *Article) ToResponse() response.ArticleListResponse {
	return response.ArticleListResponse{
		ID:          ar.ID,
		DoctorID:    ar.DoctorID,
		Title:       ar.Title,
		Content:     ar.Content,
		ImageUrl:    ar.ImageUrl,
		Date:        ar.Date,
		ViewCount:   ar.ViewCount,
		IsLiked:     ar.IsLiked,
		ReadingTime: ar.ReadingTime,
		Doctor: response.DoctorInfoResponse{
			ID:   ar.Doctor.ID,
			Name: ar.Doctor.Name,
		},
	}
}
