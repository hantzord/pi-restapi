package response

type ChatMessageResponse struct {
	Id      uint   `json:"id"`
	Message string `json:"message"`
	Role    string `json:"role"`
	Date    string `json:"date"`
}