package response

type ConsultationNotesCreateResponse struct {
	ID              uint   `json:"id"`
	ConsultationID  uint   `json:"consultation_id"`
	MusicID         uint   `json:"music_id"`
	ForumID         uint   `json:"forum_id"`
	MainPoint       string `json:"main_point"`
	NextStep        string `json:"next_step"`
	AdditionalNote  string `json:"additional_note"`
	MoodTrackerNote string `json:"mood_tracker_note"`
}