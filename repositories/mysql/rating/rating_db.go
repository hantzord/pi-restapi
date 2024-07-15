package rating

import (
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/user"

	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	UserId   uint `gorm:"type:int;index"`
	User     user.User `gorm:"foreignKey:user_id;references:id"`
	DoctorId uint `gorm:"type:int;index"`
	Doctor   doctor.Doctor `gorm:"foreignKey:doctor_id;references:id"`
	Rate   int `gorm:"type:int"`
	Message  string `gorm:"type:text"`
}