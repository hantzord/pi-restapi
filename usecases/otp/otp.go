package otp

import (
	"capstone/constants"
	otpEntities "capstone/entities/otp"
	"capstone/utilities"
)

type OtpUseCase struct {
	otpInterface otpEntities.RepositoryInterface
}

func NewOtpUseCase(otpInterface otpEntities.RepositoryInterface) *OtpUseCase {
	return &OtpUseCase{
		otpInterface: otpInterface,
	}
}

func (otpUseCase *OtpUseCase) SendOTP(otp otpEntities.Otp) error {
	if otp.Email == "" {
		return constants.ErrEmptyInputEmailOTP
	}

	codeOTP := utilities.GenerateOTP()
	expiry := utilities.GenerateExpiryTime()

	otp.Code = codeOTP
	otp.ExpiredAt = expiry

	err := otpUseCase.otpInterface.SendOTP(otp)
	if err != nil {
		return err
	}

	err = utilities.SendEmail(otp.Email, codeOTP)
	if err != nil {
		return err
	}

	return nil
}

func (otpUseCase *OtpUseCase) SendOTPChangeEmail(otp otpEntities.Otp) error {
	if otp.Email == "" {
		return constants.ErrEmptyInputEmailOTP
	}

	codeOTP := utilities.GenerateOTP()
	expiry := utilities.GenerateExpiryTime()

	otp.Code = codeOTP
	otp.ExpiredAt = expiry

	err := otpUseCase.otpInterface.SendOTPChangeEmail(otp)
	if err != nil {
		return err
	}

	err = utilities.SendEmail(otp.Email, codeOTP)
	if err != nil {
		return err
	}

	return nil
}

func (otpUseCase *OtpUseCase) VerifyOTP(otp otpEntities.Otp) error {
	if otp.Email == "" || otp.Code == "" {
		return constants.ErrEmptyInputVerifyOTP
	}

	err := otpUseCase.otpInterface.VerifyOTP(otp)
	if err != nil {
		return err
	}

	return nil
}

func (otpUseCase *OtpUseCase) VerifyOTPRegister(otp otpEntities.Otp) error {
	if otp.Email == "" || otp.Code == "" {
		return constants.ErrEmptyInputVerifyOTP
	}

	err := otpUseCase.otpInterface.VerifyOTPRegister(otp)
	if err != nil {
		return err
	}

	return nil
}

func (otpUseCase *OtpUseCase) VerifyOTPChangeEmail(userId int, otp otpEntities.Otp) error {
	if otp.Email == "" || otp.Code == "" {
		return constants.ErrEmptyInputVerifyOTP
	}

	err := otpUseCase.otpInterface.VerifyOTPChangeEmail(userId, otp)
	if err != nil {
		return err
	}

	return nil
}
