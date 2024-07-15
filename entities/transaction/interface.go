package transaction

import (
	"capstone/controllers/transaction/response"
	"capstone/entities"
)

type TransactionRepository interface {
	Insert(transaction *Transaction) (*Transaction, error)
	FindByID(ID string) (*Transaction, error)
	FindByConsultationID(consultationID uint) (*Transaction, error)
	FindAllByUserID(metadata *entities.Metadata, userID uint, status string) (*[]Transaction, error)
	FindAllByDoctorID(metadata *entities.Metadata, doctorID uint, status string) (*[]Transaction, error)
	Update(transaction *Transaction) (*Transaction, error)
	Delete(ID string) error
	CountTransactionByDoctorID(doctorID uint) (*response.TransactionCount, error)
}

type TransactionUseCase interface {
	InsertWithBuiltInInterface(transaction *Transaction, isUsePoint bool, userID int) (*Transaction, error)
	InsertWithCustomInterface(transaction *Transaction, isUsePoint bool, userID int) (*Transaction, error)
	ConfirmedPayment(id string, transactionStatus string) (*Transaction, error)
	FindByID(ID string) (*Transaction, error)
	FindByConsultationID(consultationID uint) (*Transaction, error)
	FindAllByUserID(metadata *entities.Metadata, userID uint, status string) (*[]Transaction, error)
	FindAllByDoctorID(metadata *entities.Metadata, doctorID uint, status string) (*[]Transaction, error)
	Update(transaction *Transaction) (*Transaction, error)
	Delete(ID string) error
	CountTransactionByDoctorID(doctorID uint) (*response.TransactionCount, error)
}
