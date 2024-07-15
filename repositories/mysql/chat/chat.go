package chat

import (
	"capstone/entities"
	chatEntities "capstone/entities/chat"
	"capstone/entities/doctor"
	userEntities "capstone/entities/user"
	"capstone/repositories/mysql/consultation"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{
		db: db,
	}
}

func (c *ChatRepo) CreateChatRoom(consultationId uint) (error) {
	var newChat Chat
	newChat.ConsultationId = consultationId

	return c.db.Create(&newChat).Error
}

func (c *ChatRepo) GetAllChatByUserId(userId int, metadata entities.Metadata, status string, search string) ([]chatEntities.Chat, error) {
	var consultationDB []consultation.Consultation

	query := c.db.Where("user_id = ?", userId)

	if status != "" && status == "completed" {
		query = query.Where("status = ? OR status = ?", "rejected", "done")
	} else if status != "" && status == "process" {
		query = query.Where("status = ? OR status = ?", "pending", "incoming")
	} else if status != "" && status == "active" {
		query = query.Where("status = ?", "active")
	}

	err := query.Find(&consultationDB).Error
	if err != nil {
		return nil, err
	}

	var consultationIds []int
	for _, consultation := range consultationDB {
		consultationIds = append(consultationIds, int(consultation.ID))
	}

	var chatDB []Chat

	query = c.db.Preload("Consultation").Preload("Consultation.Doctor")
	if search != "" {
		query = query.Joins("JOIN consultations ON consultations.id = chats.consultation_id").
			Joins("JOIN doctors ON doctors.id = consultations.doctor_id").
			Where("consultations.id IN ?", consultationIds).
			Where("doctors.name LIKE ?", "%"+search+"%")
	} else {
		query = query.Where("consultation_id IN ?", consultationIds)
	}

	err = query.Order("created_at DESC").Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Find(&chatDB).Error
	if err != nil {
		return nil, err
	}

	latestMessage := make([]ChatMessage, len(chatDB))

	for i, chat := range chatDB {
		var mesageTemp ChatMessage
		err := c.db.Where("chat_id = ?", chat.ID).Order("created_at DESC").Limit(1).Find(&mesageTemp).Error
		if err != nil {
			return nil, err
		}
		latestMessage[i] = mesageTemp
	}

	chatEnts := make([]chatEntities.Chat, len(chatDB))

	for i, chat := range chatDB {
		chatEnts[i].ID = chat.ID
		chatEnts[i].LatestMessageID = latestMessage[i].ID
		chatEnts[i].LatestMessageContent = latestMessage[i].Message
		chatEnts[i].LatestMessageTime = latestMessage[i].CreatedAt.Format("2006-01-02 15:04:05")
		chatEnts[i].ConsultationEndTime = chat.Consultation.EndDate.Format("2006-01-02 15:04:05")
		chatEnts[i].Consultation.Doctor = &doctor.Doctor{
			ID:         chat.Consultation.Doctor.ID,
			Name:       chat.Consultation.Doctor.Name,
			Username:   chat.Consultation.Doctor.Username,
			ProfilePicture: chat.Consultation.Doctor.ProfilePicture,
			Specialist: chat.Consultation.Doctor.Specialist,
		}

		if chat.Consultation.Status == "pending" {
			chatEnts[i].Status = "process"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "rejected" {
			chatEnts[i].Status = "completed"
			chatEnts[i].Isrejected = true
		} else if chat.Consultation.Status == "incoming" {
			chatEnts[i].Status = "incoming"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "active" {
			chatEnts[i].Status = "active"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "done" {
			chatEnts[i].Status = "completed"
			chatEnts[i].Isrejected = false
		}
	}

	return chatEnts, nil
}

func (c *ChatRepo) GetAllChatByDoctorId(doctorId int, metadata entities.Metadata, status string, search string) ([]chatEntities.Chat, error) {
	var consultationDB []consultation.Consultation

	query := c.db.Where("doctor_id = ?", doctorId)

	if status != "" && status == "completed" {
		query = query.Where("status = ? OR status = ?", "rejected", "done")
	} else if status != "" && status == "process" {
		query = query.Where("status = ? OR status = ?", "pending", "incoming")
	} else if status != "" && status == "active" {
		query = query.Where("status = ?", "active")
	}

	err := query.Find(&consultationDB).Error
	if err != nil {
		return nil, err
	}

	var consultationIds []int
	for _, consultation := range consultationDB {
		consultationIds = append(consultationIds, int(consultation.ID))
	}

	var chatDB []Chat

	query = c.db.Preload("Consultation").Preload("Consultation.User")
	if search != "" {
		query = query.Joins("JOIN consultations ON consultations.id = chats.consultation_id").
			Joins("JOIN users ON users.id = consultations.user_id").
			Where("consultations.id IN ?", consultationIds).
			Where("users.name LIKE ?", "%"+search+"%")
	} else {
		query = query.Where("consultation_id IN ?", consultationIds)
	}

	err = query.Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Find(&chatDB).Error
	if err != nil {
		return nil, err
	}

	var latestMessage []ChatMessage

	for _, chat := range chatDB {
		var mesageTemp ChatMessage
		err := c.db.Where("chat_id = ?", chat.ID).Order("created_at DESC").Limit(1).Find(&mesageTemp).Error
		if err != nil {
			return nil, err
		}
		latestMessage = append(latestMessage, mesageTemp)
	}

	chatEnts := make([]chatEntities.Chat, len(chatDB))

	for i, chat := range chatDB {
		chatEnts[i].ID = chat.ID
		chatEnts[i].LatestMessageID = latestMessage[i].ID
		chatEnts[i].LatestMessageContent = latestMessage[i].Message
		chatEnts[i].LatestMessageTime = latestMessage[i].CreatedAt.Format("2006-01-02 15:04:05")
		chatEnts[i].Consultation.User = userEntities.User{
			Id:       chat.Consultation.User.Id,
			Name:     chat.Consultation.User.Name,
			Username: chat.Consultation.User.Username,
			ProfilePicture: chat.Consultation.User.ProfilePicture,
		}

		if chat.Consultation.Status == "pending" {
			chatEnts[i].Status = "process"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "rejected" {
			chatEnts[i].Status = "completed"
			chatEnts[i].Isrejected = true
		} else if chat.Consultation.Status == "incoming" {
			chatEnts[i].Status = "incoming"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "active" {
			chatEnts[i].Status = "active"
			chatEnts[i].Isrejected = false
		} else if chat.Consultation.Status == "done" {
			chatEnts[i].Status = "completed"
			chatEnts[i].Isrejected = false
		}
	}

	return chatEnts, nil
}

func (c *ChatRepo) SendMessage(chatMessage chatEntities.ChatMessage) (chatEntities.ChatMessage, error) {
	var chatMessageDB ChatMessage
	chatMessageDB.ID = chatMessage.ID
	chatMessageDB.ChatId = chatMessage.ChatID
	chatMessageDB.Message = chatMessage.Message
	chatMessageDB.Role = chatMessage.Role

	if err := c.db.Create(&chatMessageDB).Error; err != nil {
		return chatEntities.ChatMessage{}, err
	}

	return chatEntities.ChatMessage{
		ID: chatMessageDB.ID,
		ChatID: chatMessageDB.ChatId,
		Message: chatMessageDB.Message,
		Role: chatMessageDB.Role,
		CreatedAt: chatMessageDB.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (c *ChatRepo) GetAllMessages(chatId int, lastMessageId int, metadata entities.Metadata) ([]chatEntities.ChatMessage, error) {
	var chatMessageDB []ChatMessage

	if lastMessageId == 0 {
		if err := c.db.Where("chat_id = ?", chatId).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Find(&chatMessageDB).Error; err != nil {
			return nil, err
		}
	} else {
		if err := c.db.Where("chat_id = ?", chatId).Where("id > ?", lastMessageId).Limit(metadata.Limit).Offset((metadata.Page - 1) * metadata.Limit).Find(&chatMessageDB).Error; err != nil {
			return nil, err
		}
	}

	chatMessageEnts := make([]chatEntities.ChatMessage, len(chatMessageDB))

	for i, chatMessage := range chatMessageDB {
		chatMessageEnts[i].ID = chatMessage.ID
		chatMessageEnts[i].ChatID = chatMessage.ChatId
		chatMessageEnts[i].Message = chatMessage.Message
		chatMessageEnts[i].Role = chatMessage.Role
		chatMessageEnts[i].CreatedAt = chatMessage.CreatedAt.Format("2006-01-02 15:04:05")
	}

	return chatMessageEnts, nil
}

func (c *ChatRepo) GetConsultationIdByChatId(chatId int) (int, error) {
	var chatDB Chat
	if err := c.db.Where("id = ?", chatId).Find(&chatDB).Error; err != nil {
		return 0, err
	}
	return int(chatDB.ConsultationId), nil
}