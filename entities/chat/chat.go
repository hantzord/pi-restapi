package chat

import (
	"capstone/entities"
	"capstone/entities/consultation"
)

type Chat struct {
	ID           uint
	Consultation consultation.Consultation
	Status       string
	Isrejected   bool
	LatestMessageID uint
	LatestMessageContent string
	LatestMessageTime string
	ConsultationEndTime string
}

type ChatMessage struct {
	ID uint
	ChatID uint
	Message string
	Role string
	CreatedAt string
}

type RepositoryInterface interface {
	CreateChatRoom(consultationId uint) (error)
	GetAllChatByUserId(userId int, metadata entities.Metadata, status string, search string) ([]Chat, error)
	GetAllChatByDoctorId(doctorId int, metadata entities.Metadata, status string, search string) ([]Chat, error)
	SendMessage(chatMessage ChatMessage) (ChatMessage, error)
	GetAllMessages(chatId int, lastMessageId int, metadata entities.Metadata) ([]ChatMessage, error)
	GetConsultationIdByChatId(chatId int) (int, error)
}

type UseCaseInterface interface {
	GetAllChatByUserId(userId int, metadata entities.Metadata, status string, search string) ([]Chat, error)
	GetAllChatByDoctorId(doctorId int, metadata entities.Metadata, status string, search string) ([]Chat, error)
	SendMessage(chatMessage ChatMessage) (ChatMessage, error)
	GetAllMessages(chatId int, lastMessageId string, metadata entities.Metadata) ([]ChatMessage, error)
}