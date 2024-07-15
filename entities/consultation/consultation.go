package consultation

import (
	"capstone/controllers/consultation/response"
	"capstone/entities"
	"capstone/entities/complaint"
	"capstone/entities/doctor"
	"capstone/entities/forum"
	"capstone/entities/music"
	"capstone/entities/user"
	"time"
)

type Consultation struct {
	ID            uint
	DoctorID      uint
	Doctor        *doctor.Doctor
	UserID        uint
	User          user.User
	Complaint     complaint.Complaint
	Status        string
	PaymentStatus string
	IsAccepted    bool
	IsActive      bool
	StartDate     time.Time
	EndDate       time.Time
}

type ConsultationNotes struct {
	ID              uint
	ConsultationID  uint
	Consultation    Consultation
	MusicID         uint
	Music           music.Music
	ForumID         uint
	Forum           forum.Forum
	MainPoint       string
	NextStep        string
	AdditionalNote  string
	MoodTrackerNote string
	CreatedAt       string
}

type CountConsultation struct {
	TotalConsultation    int64
	TodayConsultation    int64
	ActiveConsultation   int64
	DoneConsultation     int64
	RejectedConsultation int64
	IncomingConsultation int64
	PendingConsultation  int64
}

type ConsultationRepository interface {
	CreateConsultation(consultation *Consultation) (*Consultation, error)
	GetConsultationByID(consultationID int) (*Consultation, error)
	GetDoctorConsultationByID(consultationID int) (*Consultation, error)
	GetAllUserConsultation(metadata *entities.Metadata, userID int) (*[]Consultation, error)
	UpdateStatusConsultation(consultation *Consultation) (*Consultation, error)
	UpdatePaymentStatusConsultation(consultationID int, status string) error
	GetAllDoctorConsultation(metadata *entities.Metadata, doctorID int) (*[]Consultation, error)
	GetConsultationByComplaintID(complaintID int) (*Consultation, error)
	CountConsultationByStatus(doctorID int, status string) (int64, error)
	CountConsultationToday(doctorID int) (int64, error)
	CountConsultationByDoctorID(doctorID int) (int64, error)
	CreateConsultationNotes(consultationNotes ConsultationNotes) (ConsultationNotes, error)
	GetConsultationNotesByID(consultationID int) (ConsultationNotes, error)
	GetAllConsultation() *[]Consultation
	GetByComplaintID(complaintID int) (*Consultation, error)
	GetDoctorConsultationByComplaint(metadata *entities.Metadata, doctorID int) (*[]Consultation, error)
	SearchConsultationByComplaintName(metadata *entities.Metadata, doctorID int, name string) (*[]Consultation, error)
}

type ConsultationUseCase interface {
	CreateConsultation(consultation *Consultation) (*Consultation, error)
	GetConsultationByID(consultationID int) (*Consultation, error)
	GetDoctorConsultationByID(consultationID int) (*Consultation, error)
	GetAllUserConsultation(metadata *entities.Metadata, userID int) (*[]Consultation, error)
	UpdateStatusConsultation(consultation *Consultation) (*Consultation, error)
	GetAllDoctorConsultation(metadata *entities.Metadata, doctorID int) (*[]Consultation, error)
	GetConsultationByComplaintID(complaintID int) (*Consultation, error)
	CountConsultationByDoctorID(doctorID int) (int64, error)
	CountConsultationToday(doctorID int) (int64, error)
	CountConsultationByStatus(doctorID int, status string) (int64, error)
	CreateConsultationNotes(consultationNotes ConsultationNotes) (ConsultationNotes, error)
	CountConsultation(doctorID int) (*CountConsultation, error)
	GetConsultationNotesByID(consultationID int) (ConsultationNotes, error)
	GetByComplaintID(complaintID int) (*Consultation, error)
	GetDoctorConsultationByComplaint(metadata *entities.Metadata, doctorID int) (*[]Consultation, error)
	SearchConsultationByComplaintName(metadata *entities.Metadata, doctorID int, name string) (*[]Consultation, error)
}

func (r *Consultation) ToUserResponse() *response.ConsultationUserResponse {
	return &response.ConsultationUserResponse{
		ID:            int(r.ID),
		Doctor:        r.Doctor.ToDoctorResponse(),
		Status:        r.Status,
		PaymentStatus: r.PaymentStatus,
		IsAccepted:    r.IsAccepted,
		IsActive:      r.IsActive,
		StartDate:     r.StartDate.Local(),
		EndDate:       r.EndDate.Local(),
	}
}

func (r *Consultation) ToDoctorResponse() *response.ConsultationDoctorResponse {
	return &response.ConsultationDoctorResponse{
		ID:            int(r.ID),
		Status:        r.Status,
		PaymentStatus: r.PaymentStatus,
		IsAccepted:    r.IsAccepted,
		IsActive:      r.IsActive,
		StartDate:     r.StartDate.Local(),
		EndDate:       r.EndDate.Local(),
		UserID:        r.UserID,
		Complaint:     r.Complaint.ToResponse(),
	}
}

func (r *CountConsultation) ToResponse() *response.ConsultationCount {
	return &response.ConsultationCount{
		TotalConsultation:    r.TotalConsultation,
		TodayConsultation:    r.TodayConsultation,
		ActiveConsultation:   r.ActiveConsultation,
		DoneConsultation:     r.DoneConsultation,
		RejectedConsultation: r.RejectedConsultation,
		IncomingConsultation: r.IncomingConsultation,
		PendingConsultation:  r.PendingConsultation,
	}
}

func ToCountConsultation(args ...int64) CountConsultation {
	return CountConsultation{
		TotalConsultation:    args[0],
		TodayConsultation:    args[1],
		ActiveConsultation:   args[2],
		DoneConsultation:     args[3],
		RejectedConsultation: args[4],
		IncomingConsultation: args[5],
		PendingConsultation:  args[6],
	}
}
