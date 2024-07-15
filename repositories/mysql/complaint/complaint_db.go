package complaint

import (
	"capstone/entities/complaint"
	"gorm.io/gorm"
)

type Complaint struct {
	gorm.Model
	Name               string `gorm:"type:varchar(255);column:name"`
	Age                int    `gorm:"type:int;column:age"`
	Gender             string `gorm:"type:Enum('pria', 'wanita')"`
	Message            string `gorm:"type:text;column:message"`
	MedicalHistory     string `gorm:"type:text;column:medical_history"`
	DoctorNotification string `gorm:"type:varchar(255);column:doctor_notification"`
}

func (r *Complaint) ToEntities() *complaint.Complaint {

	return &complaint.Complaint{
		ID:                 r.ID,
		Name:               r.Name,
		Age:                r.Age,
		Gender:             r.Gender,
		Message:            r.Message,
		MedicalHistory:     r.MedicalHistory,
		DoctorNotification: r.DoctorNotification,
	}
}

func ToComplaintModel(request *complaint.Complaint) *Complaint {
	return &Complaint{
		Name:               request.Name,
		Age:                request.Age,
		Gender:             request.Gender,
		Message:            request.Message,
		MedicalHistory:     request.MedicalHistory,
		DoctorNotification: request.DoctorNotification,
	}
}
