package response

type ForumRecommendationResponse struct {
	ForumID         uint   `json:"forum_id"`
	Name            string `json:"name"`
	ImageUrl        string `json:"image_url"`
	NumberOfMembers int    `json:"number_of_members"`
}