package otp

import "time"

type Otp struct {
	ID        uint
	Email     string
	Code      string
	ExpiredAt time.Time
}

type RepositoryInterface interface {
	SendOTP(otp Otp) error
	SendOTPChangeEmail(otp Otp) error
	VerifyOTP(otp Otp) error
	VerifyOTPRegister(otp Otp) error
	VerifyOTPChangeEmail(userId int, otp Otp) error
}

type UseCaseInterface interface {
	SendOTP(otp Otp) error
	SendOTPChangeEmail(otp Otp) error
	VerifyOTP(otp Otp) error
	VerifyOTPRegister(otp Otp) error
	VerifyOTPChangeEmail(userId int, otp Otp) error
}
