package consultation

import (
	"capstone/entities/consultation"
	"capstone/repositories/mysql/complaint"
	"capstone/repositories/mysql/doctor"
	"capstone/repositories/mysql/forum"
	"capstone/repositories/mysql/music"
	"capstone/repositories/mysql/user"
	"time"

	"gorm.io/gorm"
)

type Consultation struct {
	gorm.Model
	DoctorID      uint                `gorm:"column:doctor_id;not null"`
	Doctor        doctor.Doctor       `gorm:"foreignKey:doctor_id;references:id"`
	UserID        uint                `gorm:"column:user_id;not null"`
	User          user.User           `gorm:"foreignKey:user_id;references:id"`
	ComplaintID   uint                `gorm:"column:complaint_id;unique;default:NULL"`
	Complaint     complaint.Complaint `gorm:"foreignKey:complaint_id;references:id"`
	Status        string              `gorm:"column:status;not null;default:'pending';type:enum('pending', 'rejected', 'incoming', 'active', 'done')"`
	PaymentStatus string              `gorm:"column:payment_status;not null;type:enum('pending', 'paid', 'canceled');default:'pending'"`
	IsAccepted    bool                `gorm:"column:is_accepted"`
	IsActive      bool                `gorm:"column:is_active"`
	StartDate     time.Time           `gorm:"column:start_date;type:datetime;NULL"`
	EndDate       time.Time           `gorm:"column:end_date;type:datetime;default:NULL"`
}

type ConstultationNotes struct {
	gorm.Model
	ConsultationID  uint         `gorm:"column:consultation_id;unique;not null"`
	Consultation    Consultation `gorm:"foreignKey:consultation_id;references:id"`
	MusicID         uint         `gorm:"column:music_id;default:NULL"`
	Music           music.Music  `gorm:"foreignKey:music_id;references:id"`
	ForumID         uint         `gorm:"column:forum_id;default:NULL"`
	Forum           forum.Forum  `gorm:"foreignKey:forum_id;references:id"`
	MainPoint       string       `gorm:"column:main_point;default:NULL"`
	NextStep        string       `gorm:"column:next_step;default:NULL"`
	AdditionalNote  string       `gorm:"column:additional_note;default:NULL"`
	MoodTrackerNote string       `gorm:"column:mood_tracker_note;default:NULL"`
}

func (receiver Consultation) ToEntities() *consultation.Consultation {
	return &consultation.Consultation{
		ID:            receiver.ID,
		DoctorID:      receiver.DoctorID,
		Doctor:        receiver.Doctor.ToEntities(),
		UserID:        receiver.UserID,
		Complaint:     *receiver.Complaint.ToEntities(),
		Status:        receiver.Status,
		PaymentStatus: receiver.PaymentStatus,
		IsAccepted:    receiver.IsAccepted,
		IsActive:      receiver.IsActive,
		StartDate:     receiver.StartDate,
		EndDate:       receiver.EndDate,
	}
}

func ToConsultationModel(request *consultation.Consultation) *Consultation {
	return &Consultation{
		Model:         gorm.Model{ID: request.ID},
		DoctorID:      request.DoctorID,
		UserID:        request.UserID,
		Status:        request.Status,
		PaymentStatus: request.PaymentStatus,
		IsAccepted:    request.IsAccepted,
		IsActive:      request.IsActive,
		StartDate:     request.StartDate,
		EndDate:       request.EndDate,
	}
}
