package response

import "time"

type StoriesGetDoctorResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	ImageUrl  string    `json:"image_url"`
	ViewCount int       `json:"view_count"`
}