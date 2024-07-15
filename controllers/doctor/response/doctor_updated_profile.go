package response

type DoctorUpdateProfileResponse struct {
	ID                     uint   `json:"id"`
	Username               string `json:"username"`
	Email                  string `json:"email"`
	Name                   string `json:"name"`
	Address                string `json:"address"`
	PhoneNumber            string `json:"phone_number"`
	Gender                 string `json:"gender"`
	ProfilePicture         string `json:"profile_picture"`
	Experience             int    `json:"experience"`
	BachelorAlmamater      string `json:"bachelor_almamater"`
	BachelorGraduationYear int    `json:"bachelor_graduation_year"`
	MasterAlmamater        string `json:"master_almamater"`
	MasterGraduationYear   int    `json:"master_graduation_year"`
	PracticeLocation       string `json:"practice_location"`
	PracticeCity           string `json:"practice_city"`
	Fee                    int    `json:"fee"`
	Specialist             string `json:"specialist"`
}
