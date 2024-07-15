package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id             int    `gorm:"primaryKey:autoIncrement"`
	Name           string `gorm:"type:varchar(100)"`
	Username       string `gorm:"type:varchar(100)"`
	Email          string `gorm:"type:varchar(100)"`
	Password       string `gorm:"type:varchar(100)"`
	Address        string `gorm:"type:text"`
	Bio            string `gorm:"type:text"`
	PhoneNumber    string `gorm:"type:varchar(100)"`
	Gender         string `gorm:"type:ENUM('pria', 'wanita');default:pria"`
	Age            int    `gorm:"type:int"`
	ProfilePicture string `gorm:"type:varchar(255)"`
	IsOauth        bool   `gorm:"type:boolean;default:false"`
	Points         int    `gorm:"type:int;default:0"`
	IsActive       bool   `gorm:"type:boolean;default:false"`
	PendingEmail   string `gorm:"type:varchar(100)"`
}
