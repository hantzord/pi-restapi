package request

import (
	"capstone/entities/transaction"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	ConsultationID uint `json:"consultation_id" binding:"required" validate:"required"`
	Price          int  `json:"price" binding:"required" validate:"required"`
	Bank           string
	PaymentType    string
	UsePoint       bool `json:"use_point" form:"use_point" validator:"required"`
}

func (r TransactionRequest) ToEntities() *transaction.Transaction {
	return &transaction.Transaction{
		ID:             uuid.New(),
		ConsultationID: r.ConsultationID,
		Price:          r.Price,
		Bank:           r.Bank,
		PaymentType:    r.PaymentType,
	}
}
