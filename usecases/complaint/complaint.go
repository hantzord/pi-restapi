package complaint

import (
	"capstone/entities"
	complaintEntities "capstone/entities/complaint"
	consultationEntities "capstone/entities/consultation"
	notificationEntities "capstone/entities/notification"
	"capstone/utilities"
)

type ComplaintUseCase struct {
	complaintRepo       complaintEntities.ComplaintRepository
	notificationUseCase notificationEntities.NotificationUseCase
	consultationUseCase consultationEntities.ConsultationUseCase
}

func NewComplaintUseCase(complaintRepo complaintEntities.ComplaintRepository, notificationUseCase notificationEntities.NotificationUseCase, consultationUseCase consultationEntities.ConsultationUseCase) complaintEntities.ComplaintUseCase {
	return &ComplaintUseCase{
		complaintRepo:       complaintRepo,
		notificationUseCase: notificationUseCase,
		consultationUseCase: consultationUseCase,
	}
}

func (usecase *ComplaintUseCase) Create(complaint *complaintEntities.Complaint) (*complaintEntities.Complaint, error) {
	result, err := usecase.complaintRepo.Create(complaint)
	if err != nil {
		return nil, err
	}

	consultation, err := usecase.consultationUseCase.GetConsultationByID(int(complaint.ConsultationID))
	if err != nil {
		return nil, err
	}

	notificationContent := utilities.AddContentComplaintUserNotification(result.Name, result.Message)

	if err = usecase.notificationUseCase.CreateDoctorNotification(consultation.DoctorID, notificationContent); err != nil {
		return nil, err
	}

	return result, nil
}

func (usecase *ComplaintUseCase) GetAllByDoctorID(metadata *entities.Metadata, doctorID int) (*[]complaintEntities.Complaint, error) {
	result, err := usecase.complaintRepo.GetAllByDoctorID(metadata, doctorID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ComplaintUseCase) GetByID(complaintID int) (*complaintEntities.Complaint, error) {
	result, err := usecase.complaintRepo.GetByID(complaintID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ComplaintUseCase) SearchComplaintByPatientName(metadata *entities.Metadata, name string, doctorID uint) (*[]complaintEntities.Complaint, error) {
	result, err := usecase.complaintRepo.SearchComplaintByPatientName(metadata, name, doctorID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
