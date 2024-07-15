package request

type PostLikeRequest struct {
	PostId uint `json:"post_id" form:"post_id"`
}