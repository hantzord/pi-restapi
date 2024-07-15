package response

import (
	complaintResponse "capstone/controllers/complaint/response"
	doctorResponse "capstone/controllers/doctor/response"
	"time"
)

type ConsultationUserResponse struct {
	ID            int                            `json:"id"`
	Doctor        *doctorResponse.DoctorResponse `json:"doctor"`
	Status        string                         `json:"status"`
	PaymentStatus string                         `json:"payment_status"`
	IsAccepted    bool                           `json:"is_accepted"`
	IsActive      bool                           `json:"is_active"`
	StartDate     time.Time                      `json:"start_date"`
	EndDate       time.Time                      `json:"end_date"`
}

type ConsultationDoctorResponse struct {
	ID            int                                  `json:"id"`
	Status        string                               `json:"status"`
	PaymentStatus string                               `json:"payment_status"`
	IsAccepted    bool                                 `json:"is_accepted"`
	IsActive      bool                                 `json:"is_active"`
	StartDate     time.Time                            `json:"start_date"`
	EndDate       time.Time                            `json:"end_date"`
	UserID        uint                                 `json:"user_id"`
	Complaint     *complaintResponse.ComplaintResponse `json:"complaint"`
}
