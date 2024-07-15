package request

type StoryPostRequest struct {
	Title    string `json:"title" form:"title"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image_url" form:"image_url"`
}