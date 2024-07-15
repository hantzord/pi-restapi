package response

import "capstone/controllers/consultation/response"

type UserTransactionResponse struct {
	ID           string                            `json:"id"`
	Consultation response.ConsultationUserResponse `json:"consultation"`
	Price        int                               `json:"price"`
	PaymentType  string                            `json:"payment_type"`
	PaymentLink  string                            `json:"payment_link"`
	Bank         string                            `json:"bank"`
	Status       string                            `json:"status"`
	PointSpend   int                               `json:"point_spend"`
	CreatedAt    string                            `json:"created_at"`
	UpdatedAt    string                            `json:"updated_at"`
}

type DoctorTransactionResponse struct {
	ID           string                              `json:"id"`
	Consultation response.ConsultationDoctorResponse `json:"consultation"`
	Price        int                                 `json:"price"`
	PaymentType  string                              `json:"payment_type"`
	PaymentLink  string                              `json:"payment_link"`
	Bank         string                              `json:"bank"`
	Status       string                              `json:"status"`
	PointSpend   int                                 `json:"point_spend"`
	CreatedAt    string                              `json:"created_at"`
	UpdatedAt    string                              `json:"updated_at"`
}
