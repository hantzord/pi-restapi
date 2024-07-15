package request

type StoryEditRequest struct {
	Title    string `json:"title" form:"title"`
	Singer   string `json:"singer" form:"singer"`
	ImageUrl string `json:"image_url" form:"image_url"`
}