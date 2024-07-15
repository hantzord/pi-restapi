package transaction

import (
	"capstone/entities/transaction"
	"capstone/repositories/mysql/consultation"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID uuid.UUID `gorm:"column:id;primaryKey;type:char(100)"`
	gorm.Model
	ConsultationID uint                      `gorm:"column:consultation_id;not null;unique"`
	Consultation   consultation.Consultation `gorm:"foreignKey:consultation_id;references:id"`
	Price          int                       `gorm:"column:price;not null"`
	PaymentType    string                    `gorm:"column:payment_type;not null;type:enum('gopay','bank_transfer');default:'bank_transfer'"`
	Bank           string                    `gorm:"column:bank;not null;default:'ewallet'"`
	PaymentLink    string                    `gorm:"column:payment_link;not null"`
	Status         string                    `gorm:"column:status;not null;type:enum('pending','settlement','failed', 'deny');default:'pending'"`
	PointSpend     int                       `gorm:"column:point_spend;not null;default:0"`
}

func (receiver Transaction) ToEntities() *transaction.Transaction {
	return &transaction.Transaction{
		ID:             receiver.ID,
		ConsultationID: receiver.ConsultationID,
		Consultation:   *receiver.Consultation.ToEntities(),
		Price:          receiver.Price,
		PaymentType:    receiver.PaymentType,
		PaymentLink:    receiver.PaymentLink,
		Bank:           receiver.Bank,
		Status:         receiver.Status,
		PointSpend:     receiver.PointSpend,
		CreatedAt:      receiver.CreatedAt,
		UpdatedAt:      receiver.UpdatedAt,
	}
}

func ToTransactionModel(transaction *transaction.Transaction) *Transaction {
	return &Transaction{
		ID:             transaction.ID,
		Model:          gorm.Model{CreatedAt: transaction.CreatedAt, UpdatedAt: transaction.UpdatedAt},
		ConsultationID: transaction.ConsultationID,
		Price:          transaction.Price,
		PaymentType:    transaction.PaymentType,
		PaymentLink:    transaction.PaymentLink,
		Bank:           transaction.Bank,
		Status:         transaction.Status,
		PointSpend:     transaction.PointSpend,
	}
}
