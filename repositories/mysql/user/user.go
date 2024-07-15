package user

import (
	"capstone/constants"
	userEntities "capstone/entities/user"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (userRepo *UserRepo) Register(user *userEntities.User) (userEntities.User, int64, error) {
	userDb := User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	var counterUsername, counterEmail int64
	err := userRepo.DB.Model(&userDb).Where("username = ?", userDb.Username).Count(&counterUsername).Error
	if err != nil {
		return userEntities.User{}, 0, err
	}

	if counterUsername > 0 {
		return userEntities.User{}, 1, nil
	}

	err = userRepo.DB.Model(&userDb).Where("email = ?", userDb.Email).Count(&counterEmail).Error
	if err != nil {
		return userEntities.User{}, 0, err
	}

	if counterEmail > 0 {
		return userEntities.User{}, 2, nil
	}

	err = userRepo.DB.Create(&userDb).Error
	if err != nil {
		fmt.Println(err)
		return userEntities.User{}, 0, err
	}

	userResult := userEntities.User{
		Id:       userDb.Id,
		Username: userDb.Username,
		Email:    userDb.Email,
		Password: userDb.Password,
	}

	return userResult, 0, nil
}

func (userRepo *UserRepo) Login(user *userEntities.User) (userEntities.User, error) {
	userDb := User{
		Username: user.Username,
		Password: user.Password,
	}

	password := userDb.Password

	err := userRepo.DB.Where("Username = ?", userDb.Username).First(&userDb).Error
	if err != nil {
		return userEntities.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(password))
	if err != nil {
		return userEntities.User{}, err
	}

	userResult := userEntities.User{
		Id:             userDb.Id,
		Name:           userDb.Name,
		Username:       userDb.Username,
		Email:          userDb.Email,
		Password:       userDb.Password,
		Address:        userDb.Address,
		Bio:            userDb.Bio,
		PhoneNumber:    userDb.PhoneNumber,
		Gender:         userDb.Gender,
		Age:            userDb.Age,
		ProfilePicture: userDb.ProfilePicture,
	}

	return userResult, nil
}

func (r *UserRepo) Create(email string, picture string, name string, username string) (userEntities.User, error) {
	var userDB User
	userDB.Email = email
	userDB.ProfilePicture = picture
	userDB.Name = name
	userDB.IsOauth = true
	userDB.Username = username

	err := r.DB.Create(&userDB).Error
	if err != nil {
		return userEntities.User{}, err
	}

	var userEnt userEntities.User
	userEnt.Id = userDB.Id
	userEnt.Name = userDB.Name
	userEnt.Email = userDB.Email
	userEnt.ProfilePicture = userDB.ProfilePicture
	userEnt.IsOauth = userDB.IsOauth

	return userEnt, nil
}

func (r *UserRepo) OauthFindByEmail(email string) (userEntities.User, int, error) {
	var userDB User
	if err := r.DB.Where("email = ?", email).First(&userDB).Error; err != nil {
		return userEntities.User{}, 0, err
	}

	if !userDB.IsOauth {
		return userEntities.User{}, 1, constants.ErrEmailAlreadyExist
	}

	var userEnt userEntities.User
	userEnt.Id = userDB.Id
	userEnt.Name = userDB.Name
	userEnt.Email = userDB.Email
	userEnt.ProfilePicture = userDB.ProfilePicture
	userEnt.IsOauth = userDB.IsOauth

	return userEnt, 0, nil
}

func (r *UserRepo) GetPointsByUserId(id int) (int, error) {
	var userDB User
	if err := r.DB.Where("id = ?", id).First(&userDB).Error; err != nil {
		return 0, err
	}

	return userDB.Points, nil
}

func (r *UserRepo) ResetPassword(email string, password string) error {
	err := r.DB.Model(&User{}).Where("email = ?", email).Update("password", password).Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *UserRepo) UpdateUserProfile(user *userEntities.User) (userEntities.User, error) {
	existingUser := User{}
	if err := userRepo.DB.Where("id = ?", user.Id).First(&existingUser).Error; err != nil {
		return userEntities.User{}, err
	}

	existingUser.Name = user.Name
	existingUser.Username = user.Username
	existingUser.Address = user.Address
	existingUser.Bio = user.Bio
	existingUser.PhoneNumber = user.PhoneNumber
	existingUser.Gender = user.Gender
	existingUser.Age = user.Age
	existingUser.ProfilePicture = user.ProfilePicture

	if err := userRepo.DB.Save(&existingUser).Error; err != nil {
		return userEntities.User{}, err
	}

	updatedUser := userEntities.User{
		Id:             existingUser.Id,
		Name:           existingUser.Name,
		Username:       existingUser.Username,
		Email:          existingUser.Email,
		Address:        existingUser.Address,
		Bio:            existingUser.Bio,
		PhoneNumber:    existingUser.PhoneNumber,
		Gender:         existingUser.Gender,
		Age:            existingUser.Age,
		ProfilePicture: existingUser.ProfilePicture,
		IsOauth:        existingUser.IsOauth,
		Points:         existingUser.Points,
	}

	return updatedUser, nil
}

func (r *UserRepo) ChangePassword(userId int, oldPassword, newPassword string) error {
	var userDB User
	if err := r.DB.Where("id = ?", userId).First(&userDB).Error; err != nil {
		return err
	}

	// Compare old password with hashed password in database
	err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(oldPassword))
	if err != nil {
		return constants.ErrInvalidCredentials
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password in the database
	err = r.DB.Model(&User{}).Where("id = ?", userId).Update("password", string(hashedPassword)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) UpdatePointsByUserID(id int, points int) error {
	err := r.DB.Model(&User{}).Where("id = ?", id).Update("points", points).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) ChangeEmail(userId int, email string) error {
	err := r.DB.Model(&User{}).Where("id = ?", userId).Update("pending_email", email).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetDetailedProfile(id int) (userEntities.User, error) {
	var userDB User
	if err := r.DB.Where("id = ?", id).First(&userDB).Error; err != nil {
		return userEntities.User{}, err
	}

	return userEntities.User{
		Id:             userDB.Id,
		Name:           userDB.Name,
		Username:       userDB.Username,
		Email:          userDB.Email,
		Address:        userDB.Address,
		Bio:            userDB.Bio,
		PhoneNumber:    userDB.PhoneNumber,
		Gender:         userDB.Gender,
		Age:            userDB.Age,
		ProfilePicture: userDB.ProfilePicture,
		IsOauth:        userDB.IsOauth,
		Points:         userDB.Points,
	}, nil
}