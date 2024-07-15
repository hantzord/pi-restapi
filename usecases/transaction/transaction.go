package transaction

import (
	"capstone/constants"
	responseTransaction "capstone/controllers/transaction/response"
	"capstone/entities"
	"capstone/entities/consultation"
	doctorEntities "capstone/entities/doctor"
	midtransEntities "capstone/entities/midtrans"
	transactionEntities "capstone/entities/transaction"
	userEntities "capstone/entities/user"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	transactionRepository  transactionEntities.TransactionRepository
	midtransUseCase        midtransEntities.MidtransUseCase
	consultationRepository consultation.ConsultationRepository
	doctorRepository       doctorEntities.DoctorRepositoryInterface
	userUseCase            userEntities.UseCaseInterface
	validate               *validator.Validate
}

func NewTransactionUseCase(
	transactionRepository transactionEntities.TransactionRepository,
	midtransUseCase midtransEntities.MidtransUseCase,
	consultationRepository consultation.ConsultationRepository,
	doctorRepository doctorEntities.DoctorRepositoryInterface,
	userUseCase userEntities.UseCaseInterface,
	validate *validator.Validate,
) transactionEntities.TransactionUseCase {
	return &Transaction{
		transactionRepository:  transactionRepository,
		midtransUseCase:        midtransUseCase,
		consultationRepository: consultationRepository,
		doctorRepository:       doctorRepository,
		userUseCase:            userUseCase,
		validate:               validate,
	}
}

func (usecase *Transaction) InsertWithBuiltInInterface(transaction *transactionEntities.Transaction, isUsePoint bool, userID int) (*transactionEntities.Transaction, error) {
	var err error
	if err = usecase.validate.Struct(transaction); err != nil {
		return nil, err
	}

	if transaction.Price < 0 {
		return nil, constants.ErrInvalidPrice
	}

	// Check User Point
	var point int
	if isUsePoint {
		// Get User Point
		point, err = usecase.userUseCase.GetPointsByUserId(userID)
		if err != nil {
			return nil, err
		}

		// Handling Price If Use Point
		if point >= transaction.Price {
			transaction.PointSpend = transaction.Price
			transaction.Price = 0
		} else if point < transaction.Price && point > 0 {
			transaction.PointSpend = transaction.Price - point
			transaction.Price -= point
		}
	}

	// Get Consultation Data
	_, err = usecase.consultationRepository.GetConsultationByID(int(transaction.ConsultationID))
	if err != nil {
		return nil, err
	}

	// Insert Transaction If Price == 0
	var newTransaction, response *transactionEntities.Transaction
	if transaction.Price == 0 {
		transaction.Status = constants.Success
		var transactionResponse *transactionEntities.Transaction
		transactionResponse, err = usecase.transactionRepository.Insert(transaction)
		if err != nil {
			return nil, err
		}

		// Update User Point
		err = usecase.userUseCase.UpdateSuccessPointByUserID(userID, transaction.PointSpend)
		if err != nil {
			return nil, err
		}

		response, err = usecase.ConfirmedPayment(transactionResponse.ID.String(), constants.Success)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	newTransaction, err = usecase.midtransUseCase.GenerateSnapURL(transaction)
	if err != nil {
		return nil, err
	}

	response, err = usecase.transactionRepository.Insert(newTransaction)
	if err != nil {
		return nil, err
	}
	// Update User Point
	err = usecase.userUseCase.UpdateSuccessPointByUserID(userID, transaction.PointSpend)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (usecase *Transaction) InsertWithCustomInterface(transaction *transactionEntities.Transaction, isUsePoint bool, userID int) (*transactionEntities.Transaction, error) {
	var err error
	if err = usecase.validate.Struct(transaction); err != nil {
		return nil, err
	}

	if transaction.Price < 0 {
		return nil, constants.ErrInvalidPrice
	}

	// Check User Point
	var point int
	if isUsePoint {
		// Get User Point
		point, err = usecase.userUseCase.GetPointsByUserId(userID)
		if err != nil {
			return nil, err
		}

		// Handling Price If Use Point
		if point >= transaction.Price {
			transaction.PointSpend = transaction.Price
			transaction.Price = 0
		} else {
			transaction.PointSpend = transaction.Price - point
			transaction.Price -= point
		}
	}

	// Get Consultation Data
	_, err = usecase.consultationRepository.GetConsultationByID(int(transaction.ConsultationID))
	if err != nil {
		return nil, err
	}

	// Insert Transaction If Price == 0
	var newTransaction, response *transactionEntities.Transaction
	if transaction.Price == 0 {
		transaction.Status = constants.Success
		var transactionResponse *transactionEntities.Transaction
		transactionResponse, err = usecase.transactionRepository.Insert(transaction)
		if err != nil {
			return nil, err
		}
		// Update User Point
		err = usecase.userUseCase.UpdateSuccessPointByUserID(userID, transaction.PointSpend)
		if err != nil {
			return nil, err
		}

		response, err = usecase.ConfirmedPayment(transactionResponse.ID.String(), constants.Success)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	// Insert Transaction If Price > 0
	if transaction.PaymentType == constants.BankTransfer {
		newTransaction, err = usecase.midtransUseCase.BankTransfer(transaction)
		if err != nil {
			return nil, err
		}
	} else if transaction.PaymentType == constants.GoPay {
		newTransaction, err = usecase.midtransUseCase.EWallet(transaction)
		if err != nil {
			return nil, err
		}
	}

	response, err = usecase.transactionRepository.Insert(newTransaction)
	if err != nil {
		return nil, err
	}
	// Update User Point
	err = usecase.userUseCase.UpdateSuccessPointByUserID(userID, transaction.PointSpend)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (usecase *Transaction) FindByID(ID string) (*transactionEntities.Transaction, error) {
	newTransaction, err := usecase.transactionRepository.FindByID(ID)
	if err != nil {
		return nil, err
	}
	return newTransaction, nil
}

func (usecase *Transaction) FindByConsultationID(consultationID uint) (*transactionEntities.Transaction, error) {
	newTransaction, err := usecase.transactionRepository.FindByConsultationID(consultationID)
	if err != nil {
		return nil, err
	}
	return newTransaction, nil
}

func (usecase *Transaction) FindAllByUserID(metadata *entities.Metadata, userID uint, status string) (*[]transactionEntities.Transaction, error) {
	if !(status == constants.Success || status == constants.Pending || status == constants.Deny || status == constants.Failed) {
		status = ""
	}
	newTransaction, err := usecase.transactionRepository.FindAllByUserID(metadata, userID, status)
	if err != nil {
		return nil, err
	}
	if len(*newTransaction) == 0 {
		return nil, constants.ErrDataEmpty
	}
	return newTransaction, nil
}

func (usecase *Transaction) Update(transaction *transactionEntities.Transaction) (*transactionEntities.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (usecase *Transaction) Delete(ID string) error {
	if err := usecase.transactionRepository.Delete(ID); err != nil {
		return err
	}
	return nil
}

func (usecase *Transaction) ConfirmedPayment(id string, transactionStatus string) (*transactionEntities.Transaction, error) {
	transaction, err := usecase.transactionRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(transaction)
	if *transaction == (transactionEntities.Transaction{}) {
		return nil, constants.ErrDataNotFound
	}

	if transaction.Status == transactionStatus {
		return transaction, nil
	}

	// Update Transaction Status
	transaction.Status = transactionStatus
	transactionResponse, err := usecase.transactionRepository.Update(transaction)
	err = usecase.consultationRepository.UpdatePaymentStatusConsultation(int(transactionResponse.ConsultationID), transactionResponse.Status)
	if err != nil {
		return nil, err
	}

	// Add Balance Doctor
	doctorDB, err := usecase.doctorRepository.GetDoctorByID(int(transaction.Consultation.DoctorID))
	if err != nil {
		return nil, err
	}

	if transaction.Status == constants.Success {
		totalBalance := doctorDB.Amount + transaction.Price - constants.ServiceFee
		err = usecase.doctorRepository.UpdateAmount(transaction.Consultation.DoctorID, totalBalance)
		if err != nil {
			return nil, err
		}
	} else if transaction.Status == constants.Failed || transaction.Status == constants.Deny {
		// Return User Point
		userID := transaction.Consultation.UserID
		err = usecase.userUseCase.UpdateFailedPointByUserID(int(userID), transaction.PointSpend)
		if err != nil {
			return nil, err
		}
	}
	return transactionResponse, nil
}

func (usecase *Transaction) FindAllByDoctorID(metadata *entities.Metadata, doctorID uint, status string) (*[]transactionEntities.Transaction, error) {
	if !(status == constants.Success || status == constants.Pending || status == constants.Deny || status == constants.Failed) {
		status = ""
	}
	newTransaction, err := usecase.transactionRepository.FindAllByDoctorID(metadata, doctorID, status)
	if err != nil {
		return nil, err
	}
	if len(*newTransaction) == 0 {
		return nil, constants.ErrDataEmpty
	}
	return newTransaction, nil
}

func (usecase *Transaction) CountTransactionByDoctorID(doctorID uint) (*responseTransaction.TransactionCount, error) {
	countTransaction, err := usecase.transactionRepository.CountTransactionByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}
	return countTransaction, nil
}
