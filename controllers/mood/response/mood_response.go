package response

type MoodResponse struct {
	ID       uint             `json:"id"`
	Message  string           `json:"message"`
	Date     string           `json:"date"`
	ImageUrl string           `json:"image_url"`
	MoodType MoodTypeResponse `json:"mood_type"`
}

type MoodTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}