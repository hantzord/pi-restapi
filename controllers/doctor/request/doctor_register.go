package request

import (
	"capstone/constants"
	"capstone/entities/doctor"
	"capstone/utilities"
	"mime/multipart"
)

type DoctorRegisterRequest struct {
	Username         string                `json:"username" form:"username" validate:"required,min=4,max=50"`
	Email            string                `json:"email" form:"email" validate:"required,email"`
	Password         string                `json:"password" form:"password" validate:"required,min=4,max=50"`
	Name             string                `json:"name" form:"name"`
	Address          string                `json:"address" form:"address"`
	PhoneNumber      string                `json:"phone_number" form:"phone_number"`
	Gender           string                `json:"gender" form:"gender"`
	ProfilePicture   *multipart.FileHeader `json:"profile_picture" form:"profile_picture"`
	Experience       int                   `json:"experience" form:"experience"`
	BachelorAlmamater string                `json:"bachelor_almamater" form:"bachelor_almamater"`
	BachelorGraduationYear int                  `json:"bachelor_graduation_year" form:"bachelor_graduation_year"`
	MasterAlmamater string                `json:"master_almamater" form:"master_almamater"`
	MasterGraduationYear int                  `json:"master_graduation_year" form:"master_graduation_year"`
	PracticeLocation string                `json:"practice_location" form:"practice_location"`
	PracticeCity     string                `json:"practice_city" form:"practice_city"`
	Fee              int                   `json:"fee" form:"fee"`
	Specialist       string                `json:"specialist" form:"specialist"`
}

func (r *DoctorRegisterRequest) ToDoctorEntities() (*doctor.Doctor, error) {
	var err error
	var secureUrl string
	if r.ProfilePicture != nil {
		secureUrl, err = utilities.UploadImage(r.ProfilePicture)
		if err != nil {
			return nil, constants.ErrUploadImage
		}
	}

	return &doctor.Doctor{
		Username:         r.Username,
		Email:            r.Email,
		Password:         r.Password,
		Name:             r.Name,
		ProfilePicture:   secureUrl,
		Address:          r.Address,
		PhoneNumber:      r.PhoneNumber,
		Gender:           r.Gender,
		Experience:       r.Experience,
		BachelorAlmamater: r.BachelorAlmamater,
		BachelorGraduationYear: r.BachelorGraduationYear,
		MasterAlmamater: r.MasterAlmamater,
		MasterGraduationYear: r.MasterGraduationYear,
		PracticeLocation: r.PracticeLocation,
		PracticeCity:     r.PracticeCity,
		Fee:              r.Fee,
		Specialist:       r.Specialist,
	}, nil
}
