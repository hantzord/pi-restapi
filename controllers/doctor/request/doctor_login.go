package request

import (
	"capstone/entities/doctor"
)

type DoctorLoginRequest struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (r *DoctorLoginRequest) ToDoctorLoginEntities() *doctor.Doctor {
	return &doctor.Doctor{
		Email:    r.Email,
		Username: r.Username,
		Password: r.Password,
	}
}
