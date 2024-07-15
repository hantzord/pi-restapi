package music

import (
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/user"

	"gorm.io/gorm"
)

type Music struct {
	gorm.Model
	DoctorId 	uint `gorm:"type:int;index"`
	Doctor   	doctor.Doctor `gorm:"foreignKey:doctor_id;references:id"`
	Title       string `gorm:"type:varchar(100)"`
	Singer      string `gorm:"type:varchar(100)"`
	MusicUrl    string `gorm:"type:varchar(255)"`
	ImageUrl    string `gorm:"type:varchar(255)"`
}

type MusicLikes struct {
	gorm.Model
	MusicId uint `gorm:"type:int;index"`
	Music   Music `gorm:"foreignKey:music_id;references:id"`
	UserId  uint `gorm:"type:int;index"`
	User    user.User `gorm:"foreignKey:user_id;references:id"`
}

type MusicViews struct {
	gorm.Model
	MusicId uint `gorm:"type:int;index"`
	Music   Music `gorm:"foreignKey:music_id;references:id"`
	UserId  uint `gorm:"type:int;index"`
	User    user.User `gorm:"foreignKey:user_id;references:id"`
}