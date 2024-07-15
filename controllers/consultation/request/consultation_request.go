package request

import (
	"capstone/entities/consultation"
	"time"
)

type ConsultationRequest struct {
	DoctorID uint   `json:"doctor_id" form:"doctor_id" binding:"required" validate:"required"`
	UserID   int    `validate:"required"`
	Date     string `json:"date" form:"date" binding:"required"`
	Time     string `json:"time" form:"time" binding:"required"`
}

func (r ConsultationRequest) ToEntities(consultationDate, consultationTime time.Time) *consultation.Consultation {
	timeLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		timeLocation = time.Local
	}
	newDate := time.Date(consultationDate.Year(), consultationDate.Month(), consultationDate.Day(), consultationTime.Hour(), consultationTime.Minute(), consultationTime.Second(), consultationTime.Nanosecond(), timeLocation)
	return &consultation.Consultation{
		DoctorID:  r.DoctorID,
		UserID:    uint(r.UserID),
		StartDate: newDate,
	}
}
