package response

type RatingGetResponse struct {
	Id      uint                  `json:"id"`
	User    UserRatingGetResponse `json:"user"`
	Rate    int                   `json:"rate"`
	Message string                `json:"message"`
	Date    string                `json:"date"`
}

type UserRatingGetResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	ImageUrl string `json:"image_url"`
}