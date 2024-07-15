package response

type ConsultationNotesDetailResponse struct {
	ID              uint                      `json:"id"`
	ConsultationID  uint                      `json:"consultation_id"`
	Doctor          NotesDoctorDetailResponse `json:"doctor"`
	Music           NotesMusicDetailResponse  `json:"music"`
	Forum           NotesForumDetailResponse  `json:"forum"`
	MainPoint       string                    `json:"main_point"`
	NextStep        string                    `json:"next_step"`
	AdditionalNote  string                    `json:"additional_note"`
	MoodTrackerNote string                    `json:"mood_tracker_note"`
	CreatedAt       string                    `json:"created_at"`
}

type NotesDoctorDetailResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Specialist string `json:"specialist"`
	ImageUrl   string `json:"image_url"`
}

type NotesMusicDetailResponse struct {
	ID       uint   `json:"id"`
	ImageUrl string `json:"image_url"`
	Title    string `json:"title"`
}

type NotesForumDetailResponse struct {
	ID       uint   `json:"id"`
	ImageUrl string `json:"image_url"`
	Name     string `json:"name"`
}