package chatbot

type OpenAIRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ChatHistory struct {
	PreviousMessages string
}

type UseCaseInterface interface {
	GetReplyCS(message string, chatHistory []ChatHistory) (string, error)
	GetReplyMentalHealth(message string, chatHistory []ChatHistory) (string, error)
	GetReplyTreatment(message string, chatHistory []ChatHistory) (string, error)
}