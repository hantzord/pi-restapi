package forum

import (
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/user"

	"gorm.io/gorm"
)

type Forum struct {
	gorm.Model
	DoctorID uint `gorm:"type:int;index"`
	Doctor   doctor.Doctor `gorm:"foreignKey:doctor_id;references:id"`
	Name     string        `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	ImageUrl    string `gorm:"type:varchar(255)"`
}

type ForumMember struct {
	gorm.Model
	ForumID uint `gorm:"type:int;index"`
	Forum   Forum `gorm:"foreignKey:forum_id;references:id"`
	UserID  uint `gorm:"type:int;index"`
	User    user.User  `gorm:"foreignKey:user_id;references:id"`
}