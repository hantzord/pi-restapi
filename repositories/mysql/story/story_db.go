package story

import (
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/user"
	"time"

	"gorm.io/gorm"
)

type Story struct {
	gorm.Model
	Title     string `gorm:"type:varchar(100)"`
	Content   string `gorm:"type:text"`
	Date      time.Time
	ImageUrl  string `gorm:"type:varchar(255)"`
	DoctorId  uint `gorm:"type:int;index"`
	Doctor    doctor.Doctor `gorm:"foreignKey:doctor_id;references:id"`
}

type StoryLikes struct {
	gorm.Model
	StoryId uint `gorm:"type:int;index"`
	Story   Story `gorm:"foreignKey:story_id;references:id"`
	UserId  uint `gorm:"type:int;index"`
	User    user.User `gorm:"foreignKey:user_id;references:id"`
}

type StoryViews struct {
	gorm.Model
	StoryId uint `gorm:"type:int;index"`
	Story   Story `gorm:"foreignKey:story_id;references:id"`
	UserId  uint `gorm:"type:int;index"`
	User    user.User `gorm:"foreignKey:user_id;references:id"`
}