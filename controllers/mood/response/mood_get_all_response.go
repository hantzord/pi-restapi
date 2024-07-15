package response

type MoodGetAllResponse struct {
	ID       uint             `json:"id"`
	Date     string           `json:"date"`
	MoodType MoodTypeResponse `json:"mood_type"`
}