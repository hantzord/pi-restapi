package complaint

import (
	"capstone/controllers/complaint/response"
)

type Complaint struct {
	ID                 uint
	Name               string
	Age                int
	ConsultationID     uint
	Gender             string
	Message            string
	MedicalHistory     string
	DoctorNotification string
}

func (r *Complaint) ToResponse() *response.ComplaintResponse {
	return &response.ComplaintResponse{
		ID:             r.ID,
		Name:           r.Name,
		Age:            r.Age,
		Gender:         r.Gender,
		Message:        r.Message,
		MedicalHistory: r.MedicalHistory,
	}
}
