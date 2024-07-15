package doctor

import (
	doctorEntities "capstone/entities/doctor"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	Username         string `gorm:"type:varchar(100);unique;not null"`
	Email            string `gorm:"type:varchar(100);unique;not null"`
	Password         string `gorm:"type:varchar(100)"`
	Name             string `gorm:"type:varchar(100);not null"`
	Address          string `gorm:"type:text"`
	PhoneNumber      string `gorm:"type:varchar(100)"`
	Gender           string `gorm:"type:ENUM('pria', 'wanita');default:pria"`
	IsAvailable      bool   `gorm:"type:boolean;default:true"`
	ProfilePicture   string `gorm:"type:varchar(255)"`
	Balance          int    `gorm:"type:int;default:0"`
	Experience       int    `gorm:"type:int"`
	BachelorAlmamater string `gorm:"type:varchar(100)"`
	BachelorGraduationYear int    `gorm:"type:int"`
	MasterAlmamater        string `gorm:"type:varchar(100)"`
	MasterGraduationYear   int    `gorm:"type:int"`
	PracticeLocation string `gorm:"type:text"`
	PracticeCity     string `gorm:"type:varchar(100)"`
	Fee              int    `gorm:"type:int"`
	Specialist       string `gorm:"type:varchar(100)"`
	IsOauth          bool   `gorm:"type:boolean;default:false"`
	Amount           int    `gorm:"type:int;default:0"`
}

func (doctor *Doctor) ToEntities() *doctorEntities.Doctor {
	return &doctorEntities.Doctor{
		ID:               doctor.ID,
		Username:         doctor.Username,
		Email:            doctor.Email,
		Password:         doctor.Password,
		Name:             doctor.Name,
		Address:          doctor.Address,
		PhoneNumber:      doctor.PhoneNumber,
		Gender:           doctor.Gender,
		IsAvailable:      doctor.IsAvailable,
		ProfilePicture:   doctor.ProfilePicture,
		Balance:          doctor.Balance,
		Experience:       doctor.Experience,
		BachelorAlmamater: doctor.BachelorAlmamater,
		BachelorGraduationYear: doctor.BachelorGraduationYear,
		MasterAlmamater: doctor.MasterAlmamater,
		MasterGraduationYear: doctor.MasterGraduationYear,
		PracticeLocation: doctor.PracticeLocation,
		PracticeCity:     doctor.PracticeCity,
		Fee:              doctor.Fee,
		Specialist:       doctor.Specialist,
		Amount:           doctor.Amount,
	}
}

func ToDoctorModel(request *doctorEntities.Doctor) *Doctor {
	return &Doctor{
		Username:         request.Username,
		Email:            request.Email,
		Password:         request.Password,
		Name:             request.Name,
		Address:          request.Address,
		PhoneNumber:      request.PhoneNumber,
		Gender:           request.Gender,
		ProfilePicture:   request.ProfilePicture,
		Balance:          request.Balance,
		Experience:       request.Experience,
		BachelorAlmamater: request.BachelorAlmamater,
		BachelorGraduationYear: request.BachelorGraduationYear,
		MasterAlmamater: request.MasterAlmamater,
		MasterGraduationYear: request.MasterGraduationYear,
		PracticeLocation: request.PracticeLocation,
		PracticeCity:     request.PracticeCity,
		Specialist:       request.Specialist,
		Amount:           request.Amount,
	}
}
