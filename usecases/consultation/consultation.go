package consultation

import (
	"capstone/constants"
	"capstone/entities"
	chatEntities "capstone/entities/chat"
	consultationEntities "capstone/entities/consultation"
	doctorEntities "capstone/entities/doctor"
	notificationEntities "capstone/entities/notification"
	transactionEntities "capstone/entities/transaction"
	userEntities "capstone/entities/user"
	"capstone/utilities"
	"github.com/go-playground/validator/v10"
	"time"
)

type ConsultationUseCase struct {
	consultationRepo      consultationEntities.ConsultationRepository
	transactionRepository transactionEntities.TransactionRepository
	doctorRepository      doctorEntities.DoctorRepositoryInterface
	userUseCase           userEntities.UseCaseInterface
	notificationUseCase   notificationEntities.NotificationUseCase
	validate              *validator.Validate
	chatRepo              chatEntities.RepositoryInterface
}

func NewConsultationUseCase(
	consultationRepo consultationEntities.ConsultationRepository,
	transactionRepository transactionEntities.TransactionRepository,
	userUseCase userEntities.UseCaseInterface,
	doctorRepository doctorEntities.DoctorRepositoryInterface,
	notificationUseCase notificationEntities.NotificationUseCase,
	validate *validator.Validate,
	chatRepo chatEntities.RepositoryInterface) consultationEntities.ConsultationUseCase {
	return &ConsultationUseCase{
		consultationRepo:      consultationRepo,
		transactionRepository: transactionRepository,
		userUseCase:           userUseCase,
		doctorRepository:      doctorRepository,
		notificationUseCase:   notificationUseCase,
		validate:              validate,
		chatRepo:              chatRepo,
	}
}

func (usecase *ConsultationUseCase) CreateConsultation(consultation *consultationEntities.Consultation) (*consultationEntities.Consultation, error) {
	if err := usecase.validate.Struct(consultation); err != nil {
		return nil, constants.ErrDataEmpty
	}
	result, err := usecase.consultationRepo.CreateConsultation(consultation)
	if err != nil {
		return nil, err
	}

	err = usecase.chatRepo.CreateChatRoom(result.ID)
	if err != nil {
		return nil, err
	}

	contentNotification := utilities.AddContentConsultationUserNotification(result.User.Name, result.Status)
	err = usecase.notificationUseCase.CreateUserNotification(int(result.UserID), contentNotification)
	return result, nil
}

func (usecase *ConsultationUseCase) GetConsultationByID(consultationID int) (*consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetConsultationByID(consultationID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) GetAllUserConsultation(metadata *entities.Metadata, userID int) (*[]consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetAllUserConsultation(metadata, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) UpdateStatusConsultation(consultation *consultationEntities.Consultation) (*consultationEntities.Consultation, error) {
	if consultation.Status == constants.DONE {
		jakartaTime, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			return nil, constants.ErrDataNotFound
		}
		consultation.EndDate = time.Now().In(jakartaTime)
	}
	result, err := usecase.consultationRepo.UpdateStatusConsultation(consultation)
	if err != nil {
		return nil, err
	}

	if consultation.Status == constants.REJECTED {
		transaction := new(transactionEntities.Transaction)
		transaction, err = usecase.transactionRepository.FindByConsultationID(consultation.ID)
		if err != nil {
			return nil, constants.ErrConsultationAlreadyRejected
		}

		err = usecase.userUseCase.UpdateFailedPointByUserID(int(transaction.Consultation.UserID), transaction.Price)
		if err != nil {
			return nil, err
		}

		doctorResponse, err := usecase.doctorRepository.GetDoctorByID(int(consultation.DoctorID))
		if err != nil {
			return nil, err
		}

		amountDoctor := doctorResponse.Amount - transaction.Price + constants.ServiceFee
		err = usecase.doctorRepository.UpdateAmount(consultation.DoctorID, amountDoctor)
		if err != nil {
			return nil, err
		}
	}

	contentNotification := utilities.AddContentConsultationUserNotification(result.Doctor.Name, result.Status)
	if err = usecase.notificationUseCase.CreateUserNotification(int(result.UserID), contentNotification); err != nil {
		return nil, err
	}

	return result, nil
}

func (usecase *ConsultationUseCase) GetAllDoctorConsultation(metadata *entities.Metadata, doctorID int) (*[]consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetAllDoctorConsultation(metadata, doctorID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) CountConsultationByDoctorID(doctorID int) (int64, error) {
	result, err := usecase.consultationRepo.CountConsultationByDoctorID(doctorID)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (usecase *ConsultationUseCase) CountConsultationToday(doctorID int) (int64, error) {
	result, err := usecase.consultationRepo.CountConsultationToday(doctorID)
	if err != nil {
		return 0, nil
	}
	return result, nil
}

func (usecase *ConsultationUseCase) CountConsultationByStatus(doctorID int, status string) (int64, error) {
	result, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, status)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) CreateConsultationNotes(consultationNotes consultationEntities.ConsultationNotes) (consultationEntities.ConsultationNotes, error) {
	if consultationNotes.ConsultationID == 0 {
		return consultationEntities.ConsultationNotes{}, constants.ErrInvalidConsultationID
	}

	result, err := usecase.consultationRepo.CreateConsultationNotes(consultationNotes)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (usecase *ConsultationUseCase) GetConsultationNotesByID(chatID int) (consultationEntities.ConsultationNotes, error) {
	consultationID, err := usecase.chatRepo.GetConsultationIdByChatId(chatID)
	if err != nil {
		return consultationEntities.ConsultationNotes{}, err
	}

	result, err := usecase.consultationRepo.GetConsultationNotesByID(consultationID)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) GetConsultationByComplaintID(complaintID int) (*consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetConsultationByComplaintID(complaintID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) CountConsultation(doctorID int) (*consultationEntities.CountConsultation, error) {
	totalConsultation, err := usecase.consultationRepo.CountConsultationByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}

	totalConsultationToday, err := usecase.consultationRepo.CountConsultationToday(doctorID)
	if err != nil {
		return nil, err
	}

	totalConsultationActive, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, constants.ACTIVE)
	if err != nil {
		return nil, err
	}

	totalConsultationDone, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, constants.DONE)
	if err != nil {
		return nil, err
	}

	totalConsultationRejected, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, constants.REJECTED)
	if err != nil {
		return nil, err
	}

	totalConsultationIncoming, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, constants.INCOMING)
	if err != nil {
		return nil, err
	}

	totalConsultationPending, err := usecase.consultationRepo.CountConsultationByStatus(doctorID, constants.PENDING)
	if err != nil {
		return nil, err
	}

	countConsultation := consultationEntities.ToCountConsultation(totalConsultation, totalConsultationToday, totalConsultationActive, totalConsultationDone, totalConsultationRejected, totalConsultationIncoming, totalConsultationPending)

	return &countConsultation, nil
}

func (usecase *ConsultationUseCase) GetDoctorConsultationByID(consultationID int) (*consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetDoctorConsultationByID(consultationID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) GetByComplaintID(complaintID int) (*consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetByComplaintID(complaintID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) GetDoctorConsultationByComplaint(metadata *entities.Metadata, doctorID int) (*[]consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.GetDoctorConsultationByComplaint(metadata, doctorID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (usecase *ConsultationUseCase) SearchConsultationByComplaintName(metadata *entities.Metadata, doctorID int, name string) (*[]consultationEntities.Consultation, error) {
	result, err := usecase.consultationRepo.SearchConsultationByComplaintName(metadata, doctorID, name)
	if err != nil {
		return nil, err
	}
	return result, nil
}
