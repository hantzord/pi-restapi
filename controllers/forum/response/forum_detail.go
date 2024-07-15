package response

type ForumDetailResponse struct {
	ForumID     uint   `json:"forum_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}