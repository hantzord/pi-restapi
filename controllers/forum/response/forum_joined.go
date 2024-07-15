package response

type ForumJoinedResponse struct {
	ForumID         uint         `json:"forum_id"`
	Name            string       `json:"name"`
	ImageUrl        string       `json:"image_url"`
	NumberOfMembers int          `json:"number_of_members"`
	User            []UserJoined `json:"user"`
}

type UserJoined struct {
	UserID         uint   `json:"user_id"`
	ProfilePicture string `json:"profile_picture"`
}