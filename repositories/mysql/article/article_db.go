package article

import (
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/user"
	"time"

	"gorm.io/gorm"

	articleEntities "capstone/entities/article"
)

type Article struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null"`
	Content     string `gorm:"type:text"`
	Date        time.Time
	ImageUrl    string        `gorm:"type:varchar(255)"`
	DoctorID    uint          `gorm:"type:int;not null"`
	Doctor      doctor.Doctor `gorm:"foreignKey:doctor_id;references:id"`
	ReadingTime int           `gorm:"type:int"`
}

type ArticleLikes struct {
	gorm.Model
	ArticleID uint      `gorm:"type:int;index"`
	Article   Article   `gorm:"foreignKey:article_id;references:id"`
	UserId    uint      `gorm:"type:int;index"`
	User      user.User `gorm:"foreignKey:user_id;references:id"`
}

type ArticleViews struct {
	gorm.Model
	ArticleID uint      `gorm:"type:int;index"`
	Article   Article   `gorm:"foreignKey:article_id;references:id"`
	UserId    uint      `gorm:"type:int;index"`
	User      user.User `gorm:"foreignKey:user_id;references:id"`
}

func (article *Article) ToEntities() *articleEntities.Article {
	return &articleEntities.Article{
		ID:       article.ID,
		DoctorID: article.DoctorID,
		Title:    article.Title,
		Content:  article.Content,
		ImageUrl: article.ImageUrl,
	}
}

func ToArticleModel(request *articleEntities.Article) *Article {
	return &Article{
		DoctorID: request.DoctorID,
		Title:    request.Title,
		Content:  request.Content,
		ImageUrl: request.ImageUrl,
	}
}
