package chat

import (
	"capstone/controllers/chat/request"
	"capstone/controllers/chat/response"
	chatEntities "capstone/entities/chat"
	"capstone/utilities"
	"capstone/utilities/base"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ChatController struct {
	chatUseCase chatEntities.UseCaseInterface
}

func NewChatController(chatUseCase chatEntities.UseCaseInterface) *ChatController {
	return &ChatController{
		chatUseCase: chatUseCase,
	}
}

func (chatController *ChatController) GetAllChatByUserId(c echo.Context) (error) {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	status := c.QueryParam("status")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(page, limit)

	token := c.Request().Header.Get("Authorization")
	userId, _ := utilities.GetUserIdFromToken(token)

	chats, err := chatController.chatUseCase.GetAllChatByUserId(userId, *metadata, status, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := make([]response.ChatUserResponse, len(chats))

	for i, chat := range chats {
		resp[i] = response.ChatUserResponse{
			Id:            chat.ID,
			Status:        chat.Status,
			Isrejected:    chat.Isrejected,
			EndTime:       chat.ConsultationEndTime,
			LatestMessage: response.LatestMessage{
				Id:      chat.LatestMessageID,
				Message: chat.LatestMessageContent,
				Date:    chat.LatestMessageTime,
			},
			Doctor: response.DoctorChatResponse{
				Id:         chat.Consultation.Doctor.ID,
				Name:       chat.Consultation.Doctor.Name,
				Username:   chat.Consultation.Doctor.Username,
				ImageUrl:   chat.Consultation.Doctor.ProfilePicture,
				Specialist: chat.Consultation.Doctor.Specialist,
			},
		}
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success get all chat", resp))
}

func (chatController *ChatController) GetAllChatByDoctorId(c echo.Context) (error) {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	status := c.QueryParam("status")

	search := c.QueryParam("search")

	metadata := utilities.GetMetadata(page, limit)

	token := c.Request().Header.Get("Authorization")
	doctorId, _ := utilities.GetUserIdFromToken(token)

	chats, err := chatController.chatUseCase.GetAllChatByDoctorId(doctorId, *metadata, status, search)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := make([]response.ChatDoctorResponse, len(chats))

	for i, chat := range chats {
		resp[i] = response.ChatDoctorResponse{
			Id:            chat.ID,
			Status:        chat.Status,
			LatestMessage: response.LatestMessage{
				Id:      chat.LatestMessageID,
				Message: chat.LatestMessageContent,
				Date:    chat.LatestMessageTime,
			},
			User: response.UserChatResponse{
				Id:       uint(chat.Consultation.User.Id),
				Name:     chat.Consultation.User.Name,
				Username: chat.Consultation.User.Username,
				ImageUrl: chat.Consultation.User.ProfilePicture,
			},
		}
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success get all chat", resp))
}

func (chatController *ChatController) SendMessage(c echo.Context) (error) {
	var req request.ChatSendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	var chatMessageEnt chatEntities.ChatMessage
	chatMessageEnt.ChatID = req.ChatId
	chatMessageEnt.Message = req.Message
	chatMessageEnt.Role = "user"

	result, err := chatController.chatUseCase.SendMessage(chatMessageEnt)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ChatMessageResponse
	resp.Id = result.ID
	resp.Message = result.Message
	resp.Role = result.Role
	resp.Date = result.CreatedAt

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success send message", resp))
}

func (chatController *ChatController) SendMessageDoctor(c echo.Context) (error) {
	var req request.ChatSendRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, base.NewErrorResponse(err.Error()))
	}

	var chatMessageEnt chatEntities.ChatMessage
	chatMessageEnt.ChatID = req.ChatId
	chatMessageEnt.Message = req.Message
	chatMessageEnt.Role = "doctor"

	result, err := chatController.chatUseCase.SendMessage(chatMessageEnt)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	var resp response.ChatMessageResponse
	resp.Id = result.ID
	resp.Message = result.Message
	resp.Role = result.Role
	resp.Date = result.CreatedAt

	return c.JSON(http.StatusCreated, base.NewSuccessResponse("Success send message", resp))
}

func (chatController *ChatController) GetAllMessages(c echo.Context) (error) {
	id := c.Param("chatId")
	chatId, _ := strconv.Atoi(id)

	Page := c.QueryParam("page")
	Limit := c.QueryParam("limit")
	LastMessageId := c.QueryParam("last_message_id")

	metadata := utilities.GetMetadata(Page, Limit)

	chats, err := chatController.chatUseCase.GetAllMessages(chatId, LastMessageId, *metadata)
	if err != nil {
		return c.JSON(base.ConvertResponseCode(err), base.NewErrorResponse(err.Error()))
	}

	resp := make([]response.ChatMessageResponse, len(chats))
	for i, chat := range chats {
		resp[i] = response.ChatMessageResponse{
			Id:      chat.ID,
			Message: chat.Message,
			Role:    chat.Role,
			Date:    chat.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, base.NewSuccessResponse("Success get all messages", resp))
}
