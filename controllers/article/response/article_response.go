package response

import "time"

type ArticleCreatedResponse struct {
	ID          uint               `json:"id"`
	DoctorID    uint               `json:"doctor_id"`
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	ImageUrl    string             `json:"image_url"`
	Date        time.Time          `json:"date"`
	// ViewCount   int                `json:"view_count"`
	IsLiked     bool               `json:"is_liked"`
	ReadingTime int                `json:"reading_time"`
	Doctor      DoctorInfoResponse `json:"doctor"`
}

type DoctorInfoResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ArticleListResponse struct {
	ID          uint      `json:"id"`
	DoctorID    uint      `json:"doctor_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ImageUrl    string    `json:"image_url"`
	Date        time.Time `json:"date"`
	ViewCount   int       `json:"view_count"`
	IsLiked     bool      `json:"is_liked"`
	ReadingTime int       `json:"reading_time"`
	Doctor      DoctorInfoResponse
}

type ArticleGetDoctorResponse struct {
	ID          uint      `json:"id"`
	DoctorID    uint      `json:"doctor_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	ImageUrl    string    `json:"image_url"`
	ViewCount   int       `json:"view_count"`
	ReadingTime int       `json:"reading_time"`
}

type ArticleCounter struct {
	Count int `json:"count"`
}

type ArticleEditResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	ImageUrl    string    `json:"image_url"`
	ViewCount   int       `json:"view_count"`
	ReadingTime int       `json:"reading_time"`
}
