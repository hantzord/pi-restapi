package transaction

import (
	"capstone/controllers/transaction/response"
	"capstone/entities/consultation"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID             uuid.UUID
	ConsultationID uint `validate:"required"`
	Consultation   consultation.Consultation
	Price          int
	PointSpend     int
	Status         string
	PaymentType    string
	PaymentLink    string
	Bank           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r Transaction) ToUserResponse() *response.UserTransactionResponse {
	return &response.UserTransactionResponse{
		ID:           r.ID.String(),
		Consultation: *r.Consultation.ToUserResponse(),
		Price:        r.Price,
		PaymentType:  r.PaymentType,
		PaymentLink:  r.PaymentLink,
		Bank:         r.Bank,
		Status:       r.Status,
		PointSpend:   r.PointSpend,
		CreatedAt:    r.CreatedAt.String(),
		UpdatedAt:    r.UpdatedAt.String(),
	}
}

func (r Transaction) ToDoctorResponse() *response.DoctorTransactionResponse {
	return &response.DoctorTransactionResponse{
		ID:           r.ID.String(),
		Consultation: *r.Consultation.ToDoctorResponse(),
		Price:        r.Price,
		PaymentType:  r.PaymentType,
		PaymentLink:  r.PaymentLink,
		Bank:         r.Bank,
		Status:       r.Status,
		PointSpend:   r.PointSpend,
		CreatedAt:    r.CreatedAt.String(),
		UpdatedAt:    r.UpdatedAt.String(),
	}
}
