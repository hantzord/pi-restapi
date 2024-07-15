package mysql

import (
	"capstone/repositories/mysql/article"
	"capstone/repositories/mysql/chat"
	"capstone/repositories/mysql/complaint"
	"capstone/repositories/mysql/consultation"
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/forum"
	"capstone/repositories/mysql/mood"
	"capstone/repositories/mysql/music"
	"capstone/repositories/mysql/notification"
	"capstone/repositories/mysql/otp"
	"capstone/repositories/mysql/post"
	"capstone/repositories/mysql/rating"
	"capstone/repositories/mysql/story"
	"capstone/repositories/mysql/transaction"
	"capstone/repositories/mysql/user"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBName string
	DBUser string
	DBPass string
	DBHost string
	DBPort string
}

func ConnectDB(config Config) *gorm.DB {
	dbportint, _ := strconv.Atoi(config.DBPort)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		dbportint,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	InitMigrate(db)
	return db
}

func InitMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(user.User{}, doctor.Doctor{}, consultation.Consultation{}, story.Story{}, story.StoryLikes{}, complaint.Complaint{}, transaction.Transaction{}, music.Music{}, music.MusicLikes{}, rating.Rating{}, mood.Mood{}, mood.MoodType{}, forum.Forum{}, forum.ForumMember{}, post.Post{}, post.PostLike{}, post.PostComment{}, article.Article{}, article.ArticleLikes{}, consultation.ConstultationNotes{}, chat.Chat{}, chat.ChatMessage{}, otp.Otp{}, music.MusicViews{}, story.StoryViews{}, article.ArticleViews{}, notification.DoctorNotification{}, notification.UserNotification{}); err != nil {
		log.Println("Error migrating user table")
	}
}
