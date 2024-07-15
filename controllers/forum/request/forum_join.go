package request

type ForumJoinRequest struct {
	ForumID uint `json:"forum_id" form:"forum_id"`
}