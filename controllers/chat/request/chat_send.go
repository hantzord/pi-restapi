package request

type ChatSendRequest struct {
	ChatId  uint   `json:"chat_id" form:"chat_id"`
	Message string `json:"message" form:"message"`
}