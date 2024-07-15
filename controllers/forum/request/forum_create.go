package request

type ForumCreateRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Image       string `json:"image" form:"image"`
}