package otp

import (
	"time"

	"gorm.io/gorm"
)

type Otp struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);not null"`
	Code  string `gorm:"type:varchar(255);not null"`
	ExpiredAt time.Time
}