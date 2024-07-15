package complaint

import "capstone/entities"

type ComplaintRepository interface {
	Create(complaint *Complaint) (*Complaint, error)
	GetAllByDoctorID(metadata *entities.Metadata, doctorID int) (*[]Complaint, error)
	GetByID(complaintID int) (*Complaint, error)
	SearchComplaintByPatientName(metadata *entities.Metadata, name string, doctorID uint) (*[]Complaint, error)
}

type ComplaintUseCase interface {
	Create(complaint *Complaint) (*Complaint, error)
	GetAllByDoctorID(metadata *entities.Metadata, doctorID int) (*[]Complaint, error)
	GetByID(complaintID int) (*Complaint, error)
	SearchComplaintByPatientName(metadata *entities.Metadata, name string, doctorID uint) (*[]Complaint, error)
}
