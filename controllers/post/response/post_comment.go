package response

type PostCommentResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	PostId    uint   `json:"post_id"`
	CreatedAt string `json:"created_at"`
	User      UserPostCommentResponse
}

type UserPostCommentResponse struct {
	Id             uint   `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}