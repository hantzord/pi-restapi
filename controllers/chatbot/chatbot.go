package chatbot

import (
	"capstone/entities/chatbot"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type ChatbotController struct{
	chatbotUseCases chatbot.UseCaseInterface
}

func NewChatbotController(chatbotUseCases chatbot.UseCaseInterface) *ChatbotController {
	return &ChatbotController{
		chatbotUseCases: chatbotUseCases,
	}
}

func (chatbotController *ChatbotController) ChatbotCS(c echo.Context) error {
	tokenStr := c.QueryParam("token")
    if tokenStr == "" {
        return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
    }

    // Verifikasi token
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("SECRET_JWT")), nil
    })

    if err != nil || !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
    }

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil) //berfungsi untuk melakukan upgrade dari http ke websocket
	if err != nil {
		fmt.Println("WebSocket Upgrade Error: ", err)
		return err
	}
	defer ws.Close() //berfungsi untuk otomatis close koneksi websocket apabila fungsi chatbot telah selesai

	var chatHistory []chatbot.ChatHistory

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("ReadMessage Error: ", err)
			break
		}

		aiResponse, err := chatbotController.chatbotUseCases.GetReplyCS(string(msg), chatHistory)
		if err != nil {
			fmt.Println("GetReply Error: ", err)
			break
		}

		chatHistory = append(chatHistory, chatbot.ChatHistory{PreviousMessages: string(msg)})

		err = ws.WriteMessage(websocket.TextMessage, []byte(aiResponse))
		if err != nil {
			fmt.Println("WriteMessage Error: ", err)
			break
		}
	}

	return nil
}

func (chatbotController *ChatbotController) ChatbotMentalHealth(c echo.Context) error {
	tokenStr := c.QueryParam("token")
    if tokenStr == "" {
        return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
    }

    // Verifikasi token
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("SECRET_JWT")), nil
    })

    if err != nil || !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
    }

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil) //berfungsi untuk melakukan upgrade dari http ke websocket
	if err != nil {
		fmt.Println("WebSocket Upgrade Error: ", err)
		return err
	}
	defer ws.Close() //berfungsi untuk otomatis close koneksi websocket apabila fungsi chatbot telah selesai

	var chatHistory []chatbot.ChatHistory

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("ReadMessage Error: ", err)
			break
		}

		aiResponse, err := chatbotController.chatbotUseCases.GetReplyMentalHealth(string(msg), chatHistory)
		if err != nil {
			fmt.Println("GetReply Error: ", err)
			break
		}

		chatHistory = append(chatHistory, chatbot.ChatHistory{PreviousMessages: string(msg)})

		err = ws.WriteMessage(websocket.TextMessage, []byte(aiResponse))
		if err != nil {
			fmt.Println("WriteMessage Error: ", err)
			break
		}
	}

	return nil
}

func (chatbotController *ChatbotController) ChatbotTreatment(c echo.Context) error {
	tokenStr := c.QueryParam("token")
    if tokenStr == "" {
        return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
    }

    // Verifikasi token
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("SECRET_JWT")), nil
    })

    if err != nil || !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
    }

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil) //berfungsi untuk melakukan upgrade dari http ke websocket
	if err != nil {
		fmt.Println("WebSocket Upgrade Error: ", err)
		return err
	}
	defer ws.Close() //berfungsi untuk otomatis close koneksi websocket apabila fungsi chatbot telah selesai

	var chatHistory []chatbot.ChatHistory

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("ReadMessage Error: ", err)
			break
		}

		aiResponse, err := chatbotController.chatbotUseCases.GetReplyTreatment(string(msg), chatHistory)
		if err != nil {
			fmt.Println("GetReply Error: ", err)
			break
		}

		chatHistory = append(chatHistory, chatbot.ChatHistory{PreviousMessages: string(msg)})

		err = ws.WriteMessage(websocket.TextMessage, []byte(aiResponse))
		if err != nil {
			fmt.Println("WriteMessage Error: ", err)
			break
		}
	}

	return nil
}