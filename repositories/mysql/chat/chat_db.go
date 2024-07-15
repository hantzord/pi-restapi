package chat

import (
	"capstone/repositories/mysql/consultation"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ConsultationId uint 					 `gorm:"column:consultation_id;not null"`
	Consultation   consultation.Consultation `gorm:"foreignKey:consultation_id;references:id"`
}

type ChatMessage struct {
	gorm.Model
	ChatId 		uint 	`gorm:"column:chat_id;not null"`
	Chat   		Chat 	`gorm:"foreignKey:chat_id;references:id"`
	Message 	string  `gorm:"column:message;not null"`
	Role 		string 	`gorm:"column:role;not null;type:enum('user', 'doctor');default:'user'"`
}