package request

type MoodCreate struct {
	MoodTypeId uint   `json:"mood_type_id" form:"mood_type_id"`
	Message    string `json:"message" form:"message"`
	Date       string `json:"date" form:"date"`
	ImageUrl   string `json:"image" form:"image"`
}