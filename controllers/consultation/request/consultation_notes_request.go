package request

type ConsultationNotesRequest struct {
	ConsultationID  uint   `json:"consultation_id" form:"consultation_id"`
	MusicID         uint   `json:"music_id" form:"music_id"`
	ForumID         uint   `json:"forum_id" form:"forum_id"`
	MainPoint       string `json:"main_point" form:"main_point"`
	NextStep        string `json:"next_step" form:"next_step"`
	AdditionalNote  string `json:"additional_note" form:"additional_note"`
	MoodTrackerNote string `json:"mood_tracker_note" form:"mood_tracker_note"`
}