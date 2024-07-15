package mood

import (
	"capstone/repositories/mysql/user"

	"gorm.io/gorm"
)

type Mood struct {
	gorm.Model
	UserId     uint `gorm:"type:int;index"`
	User       user.User `gorm:"foreignKey:user_id;references:id"`
	MoodTypeId uint `gorm:"type:int"`
	MoodType   MoodType `gorm:"foreignKey:mood_type_id;references:id"`
	Date       string `gorm:"type:date"`
	ImageUrl   string `gorm:"type:varchar(255)"`
	Message    string `gorm:"type:text"`
}

type MoodType struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}