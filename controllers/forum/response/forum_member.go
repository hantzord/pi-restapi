package response

type ForumMemberResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}