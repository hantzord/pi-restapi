package complaint

import (
	"capstone/constants"
	"capstone/entities"
	"capstone/entities/complaint"
	"gorm.io/gorm"
)

type ComplaintRepo struct {
	db *gorm.DB
}

func NewComplaintRepo(db *gorm.DB) complaint.ComplaintRepository {
	return &ComplaintRepo{db}
}

func (repository *ComplaintRepo) Create(complaint *complaint.Complaint) (*complaint.Complaint, error) {
	complaintModel := ToComplaintModel(complaint)
	if err := repository.db.Create(&complaintModel).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}
	if err := repository.db.Table("consultations").Where("id LIKE ?", complaint.ConsultationID).Update("complaint_id", complaintModel.ID).Error; err != nil {
		return nil, err
	}
	return complaintModel.ToEntities(), nil
}

func (repository *ComplaintRepo) GetAllByDoctorID(metadata *entities.Metadata, doctorID int) (*[]complaint.Complaint, error) {
	var complaintDB []Complaint
	if err := repository.db.
		Joins("JOIN consultations ON consultations.complaint_id = complaints.id").
		Where("consultations.doctor_id = ? AND consultations.complaint_id IS NOT NULL", doctorID). // Use '=' instead of 'LIKE' for ID comparison
		Limit(metadata.Limit).
		Offset(metadata.Offset()).
		Find(&complaintDB).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var complaints []complaint.Complaint
	for _, value := range complaintDB {
		complaints = append(complaints, *value.ToEntities())
	}
	return &complaints, nil
}

func (repository *ComplaintRepo) GetByID(complaintID int) (*complaint.Complaint, error) {
	var complaintDB Complaint
	if err := repository.db.Preload("Consultation").First(&complaintDB, "complaint_id LIKE ?", complaintID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	return complaintDB.ToEntities(), nil
}

func (repository *ComplaintRepo) SearchComplaintByPatientName(metadata *entities.Metadata, name string, doctorID uint) (*[]complaint.Complaint, error) {
	var complaintDB []Complaint
	if err := repository.db.
		Joins("JOIN consultations ON consultations.complaint_id = complaints.id").
		Where("consultations.doctor_id = ?", doctorID).
		Where("complaints.name LIKE ?", "%"+name+"%").
		Limit(metadata.Limit).
		Offset(metadata.Offset()).
		Find(&complaintDB).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	var complaints []complaint.Complaint
	for _, value := range complaintDB {
		complaints = append(complaints, *value.ToEntities())
	}
	return &complaints, nil
}
