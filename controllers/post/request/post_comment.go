package request

type PostCommentRequest struct {
	PostId  uint   `json:"post_id"`
	Content string `json:"content"`
}