package response

type ChatUserResponse struct {
	Id            uint               `json:"id"`
	Status        string             `json:"status"`
	Isrejected    bool               `json:"isrejected"`
	EndTime       string             `json:"end_time"`
	LatestMessage LatestMessage      `json:"latest_message"`
	Doctor        DoctorChatResponse `json:"doctor"`
}

type LatestMessage struct {
	Id      uint   `json:"id"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type DoctorChatResponse struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	ImageUrl   string `json:"image_url"`
	Specialist string `json:"specialist"`
}