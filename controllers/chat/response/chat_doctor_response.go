package response

type ChatDoctorResponse struct {
	Id            uint             `json:"id"`
	Status        string           `json:"status"`
	LatestMessage LatestMessage    `json:"latest_message"`
	User          UserChatResponse `json:"user"`
}

type UserChatResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	ImageUrl string `json:"image_url"`
}