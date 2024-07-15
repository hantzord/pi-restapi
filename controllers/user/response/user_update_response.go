package response

type UserUpdateProfileResponse struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Address        string `json:"address"`
	Bio            string `json:"bio"`
	PhoneNumber    string `json:"phone_number"`
	Gender         string `json:"gender"`
	Age            int    `json:"age"`
	ProfilePicture string `json:"profile_picture"`
}
