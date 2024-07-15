package otp

import (
	"capstone/constants"
	otpEntities "capstone/entities/otp"
	"capstone/repositories/mysql/user"
	"time"

	"gorm.io/gorm"
)

type OtpRepo struct {
	db *gorm.DB
}

func NewOtpRepo(db *gorm.DB) *OtpRepo {
	return &OtpRepo{
		db: db,
	}
}

func (otpRepo *OtpRepo) SendOTP(otp otpEntities.Otp) error {
	var otpDB Otp
	otpDB.Email = otp.Email
	otpDB.Code = otp.Code
	otpDB.ExpiredAt = otp.ExpiredAt
	return otpRepo.db.Create(&otpDB).Error
}

func (otpRepo *OtpRepo) SendOTPChangeEmail(otp otpEntities.Otp) error {
	var otpDB Otp
	otpDB.Email = otp.Email
	otpDB.Code = otp.Code
	otpDB.ExpiredAt = otp.ExpiredAt

	var counter int64
	err := otpRepo.db.Model(&user.User{}).Where("email = ?", otp.Email).Count(&counter).Error
	if err != nil {
		return err
	}

	if counter > 0 {
		return constants.ErrEmailAlreadyExist
	}

	return otpRepo.db.Create(&otpDB).Error
}

func (otpRepo *OtpRepo) VerifyOTP(otp otpEntities.Otp) error {
	var otpDB Otp

	err := otpRepo.db.Where("email = ? AND code = ?", otp.Email, otp.Code).First(&otpDB).Error
	if err != nil {
		return constants.ErrInvalidOTP
	}

	if otpDB.ExpiredAt.Before(time.Now()) {
		return constants.ErrExpiredOTP
	}

	return nil
}

func (otpRepo *OtpRepo) VerifyOTPRegister(otp otpEntities.Otp) error {
	var otpDB Otp
	err := otpRepo.db.Where("email = ? AND code = ?", otp.Email, otp.Code).First(&otpDB).Error
	if err != nil {
		return constants.ErrInvalidOTP
	}

	if otpDB.ExpiredAt.Before(time.Now()) {
		return constants.ErrExpiredOTP
	}

	err = otpRepo.db.Model(&user.User{}).Where("email = ?", otp.Email).Update("is_active", true).Error
	if err != nil {
		return err
	}

	return nil
}

func (otpRepo *OtpRepo) VerifyOTPChangeEmail(userId int, otp otpEntities.Otp) error {
	var otpDB Otp
	err := otpRepo.db.Where("email = ? AND code = ?", otp.Email, otp.Code).First(&otpDB).Error
	if err != nil {
		return constants.ErrInvalidOTP
	}

	if otpDB.ExpiredAt.Before(time.Now()) {
		return constants.ErrExpiredOTP
	}

	err = otpRepo.db.Model(&user.User{}).Where("id = ?", userId).Update("email", otp.Email).Error
	if err != nil {
		return err
	}
	// var userDB user.User
	// err = otpRepo.db.Where("email = ?", otp.Email).First(&userDB).Error
	// if err != nil {
	// 	return constants.ErrUserNotFound
	// }

	// Update user's email to pending email and clear pending email
	// err = otpRepo.db.Model(&userDB).Updates(map[string]interface{}{
	// 	"email":         userDB.PendingEmail,
	// 	"pending_email": "",
	// }).Error
	// if err != nil {
	// 	return err
	// }

	return nil
}
