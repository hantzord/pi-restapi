package consultation

import (
	"capstone/constants"
	"capstone/entities"
	consultationEntities "capstone/entities/consultation"
	doctorEntities "capstone/entities/doctor"
	forumEntities "capstone/entities/forum"
	musicEntities "capstone/entities/music"

	"gorm.io/gorm"
)

type ConsultationRepo struct {
	db *gorm.DB
}

func NewConsultationRepo(db *gorm.DB) consultationEntities.ConsultationRepository {
	return &ConsultationRepo{
		db: db,
	}
}

func (repository *ConsultationRepo) CreateConsultation(consultation *consultationEntities.Consultation) (*consultationEntities.Consultation, error) {
	consultationRequest := ToConsultationModel(consultation)

	if err := repository.db.Create(&consultationRequest).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}
	if err := repository.db.Preload("Doctor").First(&consultationRequest, consultationRequest.ID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	return consultationRequest.ToEntities(), nil
}

func (repository *ConsultationRepo) GetConsultationByID(consultationID int) (consultation *consultationEntities.Consultation, err error) {
	var consultationDB Consultation
	if err = repository.db.Preload("User").Preload("Doctor").First(&consultationDB, consultationID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	return consultationDB.ToEntities(), nil
}

func (repository *ConsultationRepo) GetAllUserConsultation(metadata *entities.Metadata, userID int) (*[]consultationEntities.Consultation, error) {
	var consultationDB []Consultation

	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Preload("Doctor").Find(&consultationDB, "user_id LIKE ?", userID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	var consultations []consultationEntities.Consultation
	for _, consultation := range consultationDB {
		consultations = append(consultations, *consultation.ToEntities())
	}

	return &consultations, nil
}

func (repository *ConsultationRepo) UpdateStatusConsultation(consultation *consultationEntities.Consultation) (*consultationEntities.Consultation, error) {
	var consultationDB Consultation
	if err := repository.db.Preload("Complaint").Preload("Doctor").First(&consultationDB, "id LIKE ?", consultation.ID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	if consultationDB.Status == constants.REJECTED || consultationDB.Status == constants.INCOMING {
		return nil, constants.ErrConsultationAlreadyRejected
	}

	consultationDB.Status = consultation.Status

	if consultation.Status == constants.DONE {
		if err := repository.db.Model(&consultationDB).Where("id LIKE ?", consultation.ID).Update("status", consultationDB.Status).Update("end_date", consultation.EndDate).Error; err != nil {
			return nil, err
		}
	}

	if err := repository.db.Model(&consultationDB).Where("id LIKE ?", consultation.ID).Update("status", consultationDB.Status).Error; err != nil {
		return nil, err
	}
	return consultationDB.ToEntities(), nil
}

func (repository *ConsultationRepo) GetAllDoctorConsultation(metadata *entities.Metadata, doctorID int) (*[]consultationEntities.Consultation, error) {
	var consultationbDB []Consultation

	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Preload("Complaint").Find(&consultationbDB, "doctor_id LIKE ?", doctorID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	var consultations []consultationEntities.Consultation
	for _, consultation := range consultationbDB {
		consultations = append(consultations, *consultation.ToEntities())
	}

	return &consultations, nil
}

func (repository *ConsultationRepo) CountConsultationToday(doctorID int) (int64, error) {
	var count int64
	if err := repository.db.Model(&Consultation{}).Where("doctor_id LIKE ? AND DATE(start_date) = CURDATE()", doctorID).Count(&count).Error; err != nil {
		return 0, constants.ErrDataNotFound
	}
	return count, nil
}

func (repository *ConsultationRepo) CountConsultationByDoctorID(doctorID int) (int64, error) {
	var count int64
	if err := repository.db.Model(&Consultation{}).Where("doctor_id LIKE ? AND status NOT LIKE ?", doctorID, constants.REJECTED).Count(&count).Error; err != nil {
		return 0, constants.ErrDataNotFound
	}
	return count, nil
}

func (repository *ConsultationRepo) CountConsultationByStatus(doctorID int, status string) (int64, error) {
	var count int64
	if err := repository.db.Model(&Consultation{}).Where("doctor_id LIKE ? AND status LIKE ?", doctorID, status).Count(&count).Error; err != nil {
		return 0, constants.ErrDataNotFound
	}
	return count, nil
}

func (repository *ConsultationRepo) CreateConsultationNotes(consultationNotes consultationEntities.ConsultationNotes) (consultationEntities.ConsultationNotes, error) {
	var notesDB ConstultationNotes
	notesDB.ID = consultationNotes.ID
	notesDB.ConsultationID = consultationNotes.ConsultationID
	notesDB.MusicID = consultationNotes.MusicID
	notesDB.ForumID = consultationNotes.ForumID
	notesDB.MainPoint = consultationNotes.MainPoint
	notesDB.NextStep = consultationNotes.NextStep
	notesDB.AdditionalNote = consultationNotes.AdditionalNote
	notesDB.MoodTrackerNote = consultationNotes.MoodTrackerNote

	if err := repository.db.Create(&notesDB).Error; err != nil {
		return consultationEntities.ConsultationNotes{}, constants.ErrInsertDatabase
	}

	var notesEnt consultationEntities.ConsultationNotes
	notesEnt.ID = notesDB.ID
	notesEnt.ConsultationID = notesDB.ConsultationID
	notesEnt.MusicID = notesDB.MusicID
	notesEnt.ForumID = notesDB.ForumID
	notesEnt.MainPoint = notesDB.MainPoint
	notesEnt.NextStep = notesDB.NextStep
	notesEnt.AdditionalNote = notesDB.AdditionalNote
	notesEnt.MoodTrackerNote = notesDB.MoodTrackerNote

	return notesEnt, nil
}

func (repository *ConsultationRepo) GetConsultationNotesByID(consultationID int) (consultationEntities.ConsultationNotes, error) {
	var notesDB ConstultationNotes
	err := repository.db.Preload("Music").Preload("Forum").Preload("Consultation").Preload("Consultation.Doctor").Where("consultation_id = ?", consultationID).First(&notesDB).Error

	if err != nil {
		return consultationEntities.ConsultationNotes{}, constants.ErrDataNotFound
	}

	var notesEnt consultationEntities.ConsultationNotes
	notesEnt.ID = notesDB.ID
	notesEnt.CreatedAt = notesDB.CreatedAt.Format("2006-01-02 15:04:05")

	notesEnt.Consultation = consultationEntities.Consultation{
		ID: notesDB.Consultation.ID,
		Doctor: &doctorEntities.Doctor{
			ID:   notesDB.Consultation.Doctor.ID,
			Name: notesDB.Consultation.Doctor.Name,
		},
	}

	notesEnt.Forum = forumEntities.Forum{
		ID:       notesDB.Forum.ID,
		Name:     notesDB.Forum.Name,
		ImageUrl: notesDB.Forum.ImageUrl,
	}

	notesEnt.Music = musicEntities.Music{
		Id:       notesDB.Music.ID,
		Title:    notesDB.Music.Title,
		ImageUrl: notesDB.Music.ImageUrl,
	}

	notesEnt.MainPoint = notesDB.MainPoint
	notesEnt.NextStep = notesDB.NextStep
	notesEnt.AdditionalNote = notesDB.AdditionalNote
	notesEnt.MoodTrackerNote = notesDB.MoodTrackerNote

	return notesEnt, nil
}

func (repository *ConsultationRepo) GetConsultationByComplaintID(complaintID int) (*consultationEntities.Consultation, error) {
	var consultationDB Consultation
	if err := repository.db.Preload("Complaint").First(&consultationDB, "complaint_id LIKE ?", complaintID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	return consultationDB.ToEntities(), nil
}

func (repository *ConsultationRepo) UpdatePaymentStatusConsultation(consultationID int, status string) error {
	var consultationDB Consultation

	consultationDB.PaymentStatus = status

	if err := repository.db.Model(&consultationDB).Where("id LIKE ?", consultationID).Error; err != nil {
		return err
	}

	return nil
}

func (repository *ConsultationRepo) GetAllConsultation() *[]consultationEntities.Consultation {
	var consultationDB []Consultation
	if err := repository.db.Find(&consultationDB, "status NOT IN (?)", []string{constants.DONE, constants.REJECTED}).Error; err != nil {
		return nil
	}

	var consultations []consultationEntities.Consultation
	for _, consultation := range consultationDB {
		consultations = append(consultations, *consultation.ToEntities())
	}

	return &consultations
}

func (repository *ConsultationRepo) GetDoctorConsultationByID(consultationID int) (*consultationEntities.Consultation, error) {
	var consultationDB Consultation
	if err := repository.db.Preload("Complaint").First(&consultationDB, "id LIKE ?", consultationID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	return consultationDB.ToEntities(), nil
}

func (repository *ConsultationRepo) GetByComplaintID(complaintID int) (*consultationEntities.Consultation, error) {
	var consultationDB Consultation
	if err := repository.db.Preload("Complaint").First(&consultationDB, "complaint_id LIKE ?", complaintID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	return consultationDB.ToEntities(), nil
}

func (repository *ConsultationRepo) GetDoctorConsultationByComplaint(metadata *entities.Metadata, doctorID int) (*[]consultationEntities.Consultation, error) {
	var consultationDB []Consultation
	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Preload("Complaint").Find(&consultationDB, "doctor_id LIKE ? AND complaint_id IS NOT NULL", doctorID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	var consultations []consultationEntities.Consultation
	for _, consultation := range consultationDB {
		consultations = append(consultations, *consultation.ToEntities())
	}

	return &consultations, nil
}

func (repository *ConsultationRepo) SearchConsultationByComplaintName(metadata *entities.Metadata, doctorID int, name string) (*[]consultationEntities.Consultation, error) {
	var consultationDB []Consultation
	if err := repository.db.
		Limit(metadata.Limit).
		Offset(metadata.Offset()).
		Preload("Complaint").
		Joins("JOIN complaints ON consultations.complaint_id = complaints.id").
		Where("complaints.name LIKE ?", "%"+name+"%").
		Find(&consultationDB, "doctor_id LIKE ? AND complaint_id IS NOT NULL", doctorID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	var consultations []consultationEntities.Consultation
	for _, consultation := range consultationDB {
		consultations = append(consultations, *consultation.ToEntities())
	}

	return &consultations, nil
}
