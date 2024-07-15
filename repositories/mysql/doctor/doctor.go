package doctor

import (
	"capstone/constants"
	"capstone/entities"
	doctorEntities "capstone/entities/doctor"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DoctorRepo struct {
	db *gorm.DB
}

func NewDoctorRepo(db *gorm.DB) doctorEntities.DoctorRepositoryInterface {
	return &DoctorRepo{
		db: db,
	}
}

func (repository *DoctorRepo) Register(doctor *doctorEntities.Doctor) (*doctorEntities.Doctor, error) {

	doctorDb := ToDoctorModel(doctor)
	if err := repository.db.Model(&doctorDb).First(&doctorDb, "username = ?", doctorDb.Username).Error; err == nil {
		return nil, constants.ErrUsernameAlreadyExist
	}

	if err := repository.db.Model(&doctorDb).First(&doctorDb, "email = ?", doctorDb.Email).Error; err == nil {
		return nil, constants.ErrEmailAlreadyExist
	}

	if err := repository.db.Create(&doctorDb).Error; err != nil {
		return nil, constants.ErrInsertDatabase
	}

	doctorResult := doctorDb.ToEntities()
	return doctorResult, nil
}

func (repository *DoctorRepo) Login(doctor *doctorEntities.Doctor) (*doctorEntities.Doctor, error) {
	doctorDb := ToDoctorModel(doctor)

	doctorPassword := doctorDb.Password
	if err := repository.db.First(&doctorDb, "username LIKE ? OR email LIKE ?", doctorDb.Username, doctorDb.Email).Error; err != nil {
		return nil, constants.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(doctorDb.Password), []byte(doctorPassword)); err != nil {
		return nil, constants.ErrUserNotFound
	}

	result := doctorDb.ToEntities()

	return result, nil

}

func (repository *DoctorRepo) GetDoctorByID(doctorID int) (doctor *doctorEntities.Doctor, err error) {
	var doctorDb Doctor
	if err = repository.db.First(&doctorDb, doctorID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}

	doctor = doctorDb.ToEntities()
	return doctor, nil
}

func (repository *DoctorRepo) GetAllDoctor(metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	var doctorsDb []Doctor
	if err := repository.db.Limit(metadata.Limit).Offset((metadata.Page-1)*metadata.Limit).Find(&doctorsDb, "").Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var doctorsResponse []doctorEntities.Doctor
	for _, doctor := range doctorsDb {
		doctorsResponse = append(doctorsResponse, *doctor.ToEntities())
	}
	return &doctorsResponse, nil
}

func (repository *DoctorRepo) GetActiveDoctor(metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	var doctorsDb []Doctor
	if err := repository.db.Limit(metadata.Limit).Offset((metadata.Page-1)*metadata.Limit).Find(&doctorsDb, "is_available = ?", true).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var doctorsResponse []doctorEntities.Doctor
	for _, doctor := range doctorsDb {
		doctorsResponse = append(doctorsResponse, *doctor.ToEntities())
	}
	return &doctorsResponse, nil
}

func (r *DoctorRepo) Create(email string, picture string, name string, username string) (doctorEntities.Doctor, error) {
	var doctorDB Doctor
	doctorDB.Email = email
	doctorDB.ProfilePicture = picture
	doctorDB.Name = name
	doctorDB.IsOauth = true
	doctorDB.Username = username

	err := r.db.Create(&doctorDB).Error
	if err != nil {
		return doctorEntities.Doctor{}, err
	}

	var doctorEnt doctorEntities.Doctor
	doctorEnt.ID = doctorDB.ID
	doctorEnt.Name = doctorDB.Name
	doctorEnt.Email = doctorDB.Email
	doctorEnt.ProfilePicture = doctorDB.ProfilePicture
	doctorEnt.IsOauth = doctorDB.IsOauth

	return doctorEnt, nil
}

func (r *DoctorRepo) OauthFindByEmail(email string) (doctorEntities.Doctor, int, error) {
	var doctorDB Doctor
	if err := r.db.Where("email = ?", email).First(&doctorDB).Error; err != nil {
		return doctorEntities.Doctor{}, 0, err
	}

	if !doctorDB.IsOauth {
		return doctorEntities.Doctor{}, 1, constants.ErrEmailAlreadyExist
	}

	var doctorEnt doctorEntities.Doctor
	doctorEnt.ID = doctorDB.ID
	doctorEnt.Name = doctorDB.Name
	doctorEnt.Email = doctorDB.Email
	doctorEnt.ProfilePicture = doctorDB.ProfilePicture
	doctorEnt.IsOauth = doctorDB.IsOauth

	return doctorEnt, 0, nil
}

func (r *DoctorRepo) UpdateAmount(doctorID uint, amount int) error {
	if err := r.db.Model(&Doctor{}).Where("id = ?", doctorID).Update("amount", amount).Error; err != nil {
		return err
	}

	return nil
}

func (repository *DoctorRepo) SearchDoctor(search string, metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	var doctorsDb []Doctor
	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Find(&doctorsDb, "name LIKE ?", "%"+search+"%").Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var doctorsResponse []doctorEntities.Doctor
	for _, doctor := range doctorsDb {
		doctorsResponse = append(doctorsResponse, *doctor.ToEntities())
	}
	return &doctorsResponse, nil
}

func (r *DoctorRepo) UpdateDoctorProfile(doctor *doctorEntities.Doctor) (doctorEntities.Doctor, error) {
	existingDoctor := Doctor{}
	if err := r.db.Where("id = ?", doctor.ID).First(&existingDoctor).Error; err != nil {
		return doctorEntities.Doctor{}, err
	}

	// existingDoctor.Username = doctor.Username
	existingDoctor.Name = doctor.Name
	// existingDoctor.Address = doctor.Address
	// existingDoctor.PhoneNumber = doctor.PhoneNumber
	existingDoctor.Gender = doctor.Gender
	if doctor.ProfilePicture != "" {
		existingDoctor.ProfilePicture = doctor.ProfilePicture
	}
	// existingDoctor.Experience = doctor.Experience
	existingDoctor.BachelorAlmamater = doctor.BachelorAlmamater
	// existingDoctor.BachelorGraduationYear = doctor.BachelorGraduationYear
	existingDoctor.MasterAlmamater = doctor.MasterAlmamater
	// existingDoctor.MasterGraduationYear = doctor.MasterGraduationYear
	existingDoctor.PracticeLocation = doctor.PracticeLocation
	// existingDoctor.PracticeCity = doctor.PracticeCity
	// existingDoctor.Fee = doctor.Fee
	existingDoctor.Specialist = doctor.Specialist

	if err := r.db.Save(&existingDoctor).Error; err != nil {
		return doctorEntities.Doctor{}, err
	}

	updatedDoctor := doctorEntities.Doctor{
		ID:               existingDoctor.ID,
		Username:         existingDoctor.Username,
		Email:            existingDoctor.Email,
		Name:             existingDoctor.Name,
		Address:          existingDoctor.Address,
		PhoneNumber:      existingDoctor.PhoneNumber,
		Gender:           existingDoctor.Gender,
		ProfilePicture:   existingDoctor.ProfilePicture,
		Experience:       existingDoctor.Experience,
		BachelorAlmamater: existingDoctor.BachelorAlmamater,
		BachelorGraduationYear: existingDoctor.BachelorGraduationYear,
		MasterAlmamater: existingDoctor.MasterAlmamater,
		MasterGraduationYear: existingDoctor.MasterGraduationYear,
		PracticeLocation: existingDoctor.PracticeLocation,
		PracticeCity:     existingDoctor.PracticeCity,
		Fee:              existingDoctor.Fee,
		Specialist:       existingDoctor.Specialist,
	}

	return updatedDoctor, nil
}

func (r *DoctorRepo) GetDetailProfile(doctorID uint) (doctorEntities.Doctor, error) {
	var doctorDB Doctor
	if err := r.db.Where("id = ?", doctorID).First(&doctorDB).Error; err != nil {
		return doctorEntities.Doctor{}, err
	}

	var doctorEnt doctorEntities.Doctor
	doctorEnt.ID = doctorDB.ID
	doctorEnt.Username = doctorDB.Username
	doctorEnt.Email = doctorDB.Email
	doctorEnt.Name = doctorDB.Name
	doctorEnt.Address = doctorDB.Address
	doctorEnt.PhoneNumber = doctorDB.PhoneNumber
	doctorEnt.Gender = doctorDB.Gender
	doctorEnt.IsAvailable = doctorDB.IsAvailable
	doctorEnt.ProfilePicture = doctorDB.ProfilePicture
	doctorEnt.Balance = doctorDB.Balance
	doctorEnt.Experience = doctorDB.Experience
	doctorEnt.BachelorAlmamater = doctorDB.BachelorAlmamater
	doctorEnt.BachelorGraduationYear = doctorDB.BachelorGraduationYear
	doctorEnt.MasterAlmamater = doctorDB.MasterAlmamater
	doctorEnt.MasterGraduationYear = doctorDB.MasterGraduationYear
	doctorEnt.PracticeLocation = doctorDB.PracticeLocation
	doctorEnt.PracticeCity = doctorDB.PracticeCity
	doctorEnt.Fee = doctorDB.Fee
	doctorEnt.Specialist = doctorDB.Specialist
	doctorEnt.Amount = doctorDB.Amount

	return doctorEnt, nil
}