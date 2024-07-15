package response

type PostResponse struct {
	ID               uint             `json:"id"`
	Content          string           `json:"content"`
	ImageUrl         string           `json:"image_url"`
	NumberOfComments int              `json:"total_comments"`
	IsLiked          bool             `json:"is_liked"`
	User             UserPostResponse `json:"user"`
}

type PostCreateResponse struct {
	ID       uint             `json:"id"`
	ForumId  uint             `json:"forum_id"`
	Content  string           `json:"content"`
	ImageUrl string           `json:"image_url"`
	User     UserPostResponse `json:"user"`
}

type UserPostResponse struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}