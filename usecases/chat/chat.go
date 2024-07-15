package chat

import (
	"capstone/constants"
	"capstone/entities"
	chatEntities "capstone/entities/chat"
	"strconv"
)

type ChatUseCase struct {
	chatInterface chatEntities.RepositoryInterface
}

func NewChatUseCase(chatInterface chatEntities.RepositoryInterface) *ChatUseCase {
	return &ChatUseCase{
		chatInterface: chatInterface,
	}
}

func (chatUseCase *ChatUseCase) GetAllChatByUserId(userId int, metadata entities.Metadata, status string, search string) ([]chatEntities.Chat, error) {
	chats, err := chatUseCase.chatInterface.GetAllChatByUserId(userId, metadata, status, search)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (chatUseCase *ChatUseCase) GetAllChatByDoctorId(doctorId int, metadata entities.Metadata, status string, search string) ([]chatEntities.Chat, error) {
	chats, err := chatUseCase.chatInterface.GetAllChatByDoctorId(doctorId, metadata, status, search)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (chatUseCase *ChatUseCase) SendMessage(chatMessage chatEntities.ChatMessage) (chatEntities.ChatMessage, error) {
	if chatMessage.ChatID == 0 && chatMessage.Message == "" {
		return chatEntities.ChatMessage{}, constants.ErrEmptyChat
	}

	chat, err := chatUseCase.chatInterface.SendMessage(chatMessage)
	if err != nil {
		return chatEntities.ChatMessage{}, err
	}
	return chat, nil
}

func (chatUseCase *ChatUseCase) GetAllMessages(chatId int, lastMessageId string, metadata entities.Metadata) ([]chatEntities.ChatMessage, error) {
	var lastMessageIdInt int
	if lastMessageId == "" {
		lastMessageIdInt = 0
	} else {
		lastMessageIdInt, _ = strconv.Atoi(lastMessageId)
	}
	messages, err := chatUseCase.chatInterface.GetAllMessages(chatId, lastMessageIdInt, metadata)
	if err != nil {
		return nil, err
	}
	return messages, nil
}