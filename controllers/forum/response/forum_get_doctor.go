package response

type ForumGetDoctorResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	ImageUrl        string `json:"image_url"`
	NumberOfMembers int    `json:"number_of_members"`
}