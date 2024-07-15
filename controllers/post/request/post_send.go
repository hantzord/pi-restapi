package request

type PostSendRequest struct {
	ForumId  uint   `json:"forum_id" form:"forum_id"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image" form:"image"`
}