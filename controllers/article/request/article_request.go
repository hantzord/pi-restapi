package request

import (
	"capstone/entities/article"
)

type CreateArticleRequest struct {
	Title    string `json:"title" form:"title"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image_url" form:"image_url"`
	DoctorID uint   `json:"doctor_id" form:"doctor_id"`
}

func (r *CreateArticleRequest) ToArticleEntities() *article.Article {
	return &article.Article{
		Title:    r.Title,
		Content:  r.Content,
		ImageUrl: r.ImageUrl,
		DoctorID: r.DoctorID,
	}
}

type UpdateArticleRequest struct {
	Title    string `json:"title" form:"title"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image_url" form:"image_url"`
}
