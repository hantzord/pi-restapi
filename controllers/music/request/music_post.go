package request

type MusicPostRequest struct {
	Title    string `json:"title" form:"title"`
	Singer   string `json:"singer" form:"singer"`
	MusicUrl string `json:"music_url" form:"music_url"`
	ImageUrl string `json:"image_url" form:"image_url"`
}