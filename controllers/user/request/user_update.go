package request

type UpdateProfileRequest struct {
	Name           string `json:"name" form:"name"`
	Username       string `json:"username" form:"username"`
	Address        string `json:"address" form:"address"`
	Bio            string `json:"bio" form:"bio"`
	PhoneNumber    string `json:"phone_number" form:"phone_number"`
	Gender         string `json:"gender" form:"gender"`
	Age            int    `json:"age" form:"age"`
	ProfilePicture string `json:"profile_picture" form:"profile_picture"`
}
