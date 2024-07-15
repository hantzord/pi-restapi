package request

import "capstone/entities/complaint"

type ComplaintRequest struct {
	ConsultationID uint   `json:"consultation_id" form:"consultation_id" validate:"required"`
	Name           string `json:"name" form:"name" validate:"required"`
	Age            int    `json:"age" form:"age" validate:"required"`
	Gender         string `json:"gender" form:"gender" validate:"required"`
	Message        string `json:"message" form:"message" validate:"required"`
	MedicalHistory string `json:"medical_history" form:"medical_history" validate:"required"`
}

func (r *ComplaintRequest) ToEntities() *complaint.Complaint {
	return &complaint.Complaint{
		Name:           r.Name,
		Age:            r.Age,
		Gender:         r.Gender,
		Message:        r.Message,
		MedicalHistory: r.MedicalHistory,
		ConsultationID: r.ConsultationID,
	}
}
