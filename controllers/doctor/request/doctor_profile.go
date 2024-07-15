package request

type UpdateDoctorProfileRequest struct {
	// Username               string `json:"username" form:"username"`
	Name string `json:"name" form:"name"`
	// Address                string `json:"address" form:"address"`
	// PhoneNumber            string `json:"phone_number" form:"phone_number"`
	Gender         string `json:"gender" form:"gender"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
	// Experience             int    `json:"experience" form:"experience"`
	BachelorAlmamater string `json:"bachelor_almamater" form:"bachelor_almamater"`
	// BachelorGraduationYear int    `json:"bachelor_graduation_year" form:"bachelor_graduation_year"`
	MasterAlmamater string `json:"master_almamater" form:"master_almamater"`
	// MasterGraduationYear   int    `json:"master_graduation_year" form:"master_graduation_year"`
	PracticeLocation string `json:"practice_location" form:"practice_location"`
	// PracticeCity           string `json:"practice_city" form:"practice_city"`
	// Fee                    int    `json:"fee" form:"fee"`
	Specialist string `json:"specialist" form:"specialist"`
}
