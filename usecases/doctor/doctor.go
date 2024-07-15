package doctor

import (
	"capstone/constants"
	"capstone/entities"
	doctorEntities "capstone/entities/doctor"
	ratingEntities "capstone/entities/rating"
	"capstone/middlewares"
	"capstone/utilities"
	"context"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	myoauth "golang.org/x/oauth2"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type DoctorUseCase struct {
	doctorRepository doctorEntities.DoctorRepositoryInterface
	ratingRepository ratingEntities.RepositoryInterface
	oauthConfig      *myoauth.Config
	oauthConfigFB    *myoauth.Config
}

func NewDoctorUseCase(doctorRepository doctorEntities.DoctorRepositoryInterface, ratingRepository ratingEntities.RepositoryInterface ,oauthConfig *myoauth.Config, oauthConfigFB *myoauth.Config) doctorEntities.DoctorUseCaseInterface {
	return &DoctorUseCase{
		doctorRepository: doctorRepository,
		ratingRepository: ratingRepository,
		oauthConfig:      oauthConfig,
		oauthConfigFB:    oauthConfigFB,
	}
}

func (usecase *DoctorUseCase) Register(doctor *doctorEntities.Doctor) (*doctorEntities.Doctor, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, constants.ErrHashedPassword
	}

	doctor.Password = string(hashedPassword)

	doctorResult, err := usecase.doctorRepository.Register(doctor)
	if err != nil {
		return nil, err
	}
	token, err := middlewares.CreateToken(int(doctorResult.ID))
	if err != nil {
		return nil, err
	}
	doctorResult.Token = token
	return doctorResult, nil

}

func (usecase *DoctorUseCase) Login(doctor *doctorEntities.Doctor) (*doctorEntities.Doctor, error) {
	if (doctor.Email == "" && doctor.Password == "") || (doctor.Username == "" && doctor.Password == "") {
		return nil, constants.ErrEmptyInputLogin
	}
	userResult, err := usecase.doctorRepository.Login(doctor)
	if err != nil {
		return nil, err
	}
	token, err := middlewares.CreateToken(int(userResult.ID))
	if err != nil {
		return nil, err
	}
	userResult.Token = token
	return userResult, nil
}

func (usecase *DoctorUseCase) GetDoctorByID(doctorID int) (*doctorEntities.Doctor, error) {
	result, err := usecase.doctorRepository.GetDoctorByID(doctorID)
	if err != nil {
		return nil, err
	}

	ratingResult, _ := usecase.ratingRepository.GetSummaryRating(result.ID)

	result.RatingPrecentage = (ratingResult.Average / 5) * 100

	return result, nil
}

func (usecase *DoctorUseCase) GetAllDoctor(metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	result, err := usecase.doctorRepository.GetAllDoctor(metadata)
	if err != nil {
		return nil, err
	}

	ratingResult := make([]ratingEntities.Rating, len(*result))

	for i, doctor := range *result {
		ratingResult[i], _ = usecase.ratingRepository.GetSummaryRating(doctor.ID)
	}

	for i, rating := range ratingResult {
		(*result)[i].RatingPrecentage = (rating.Average / 5) * 100
	}

	return result, nil
}

func (usecase *DoctorUseCase) GetActiveDoctor(metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	result, err := usecase.doctorRepository.GetActiveDoctor(metadata)
	if err != nil {
		return nil, err
	}

	ratingResult := make([]ratingEntities.Rating, len(*result))

	for i, doctor := range *result {
		ratingResult[i], _ = usecase.ratingRepository.GetSummaryRating(doctor.ID)
	}

	for i, rating := range ratingResult {
		(*result)[i].RatingPrecentage = (rating.Average / 5) * 100
	}

	return result, nil
}

func (u *DoctorUseCase) HandleGoogleLogin() string {
	return u.oauthConfig.AuthCodeURL("state-token", myoauth.AccessTypeOffline)
}

func (u *DoctorUseCase) HandleGoogleCallback(ctx context.Context, code string) (doctorEntities.Doctor, error) {
	token, err := u.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return doctorEntities.Doctor{}, constants.ErrExcange
	}

	// Membuat layanan OAuth2
	oauth2Service, err := oauth2.NewService(ctx, option.WithTokenSource(u.oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return doctorEntities.Doctor{}, constants.ErrNewServiceGoogle
	}

	userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)
	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		return doctorEntities.Doctor{}, constants.ErrNewUserInfo
	}

	// Cek apakah pengguna sudah ada di database
	result, myCode, err := u.doctorRepository.OauthFindByEmail(userInfo.Email)
	if err != nil && myCode == 0 {

		username := utilities.GetFirstNameWithNumbers(userInfo.Name)

		newUser, err := u.doctorRepository.Create(userInfo.Email, userInfo.Picture, userInfo.Name, username)
		if err != nil {
			return doctorEntities.Doctor{}, constants.ErrInsertOAuth
		}

		tokenJWT, _ := middlewares.CreateToken(int(newUser.ID))
		newUser.Token = tokenJWT

		return newUser, nil
	}

	if err != nil && myCode == 1 {
		return doctorEntities.Doctor{}, err
	}

	tokenJWT, _ := middlewares.CreateToken(int(result.ID))
	result.Token = tokenJWT

	return result, nil
}

func (u *DoctorUseCase) HandleFacebookLogin() string {
	return u.oauthConfigFB.AuthCodeURL("state-token", myoauth.AccessTypeOffline)
}

func (u *DoctorUseCase) HandleFacebookCallback(ctx context.Context, code string) (doctorEntities.Doctor, error) {
	token, err := u.oauthConfigFB.Exchange(ctx, code)
	if err != nil {
		return doctorEntities.Doctor{}, constants.ErrExcange
	}

	client := u.oauthConfigFB.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return doctorEntities.Doctor{}, constants.ErrNewUserInfo
	}
	defer resp.Body.Close()

	var fbUser struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fbUser); err != nil {
		return doctorEntities.Doctor{}, err
	}

	result, myCode, err := u.doctorRepository.OauthFindByEmail(fbUser.Email)
	if err != nil && myCode == 0 {
		username := utilities.GetFirstNameWithNumbers(fbUser.Name)
		newUser, err := u.doctorRepository.Create(fbUser.Email, fbUser.Picture.Data.URL, fbUser.Name, username)
		if err != nil {
			return doctorEntities.Doctor{}, constants.ErrInsertOAuth
		}

		tokenJWT, _ := middlewares.CreateToken(int(newUser.ID))
		newUser.Token = tokenJWT

		return newUser, nil
	}

	if err != nil && myCode == 1 {
		return doctorEntities.Doctor{}, err
	}

	tokenJWT, _ := middlewares.CreateToken(int(result.ID))
	result.Token = tokenJWT

	return result, nil
}

func (usecase *DoctorUseCase) SearchDoctor(search string, metadata *entities.Metadata) (*[]doctorEntities.Doctor, error) {
	result, err := usecase.doctorRepository.SearchDoctor(search, metadata)
	if err != nil {
		return nil, err
	}

	ratingResult := make([]ratingEntities.Rating, len(*result))

	for i, doctor := range *result {
		ratingResult[i], _ = usecase.ratingRepository.GetSummaryRating(doctor.ID)
	}

	for i, rating := range ratingResult {
		(*result)[i].RatingPrecentage = (rating.Average / 5) * 100
	}

	return result, nil
}

func (u *DoctorUseCase) UpdateDoctorProfile(doctor *doctorEntities.Doctor) (doctorEntities.Doctor, error) {
	if doctor.ID == 0 {
		return doctorEntities.Doctor{}, constants.ErrUserNotFound
	}

	updatedDoctor, err := u.doctorRepository.UpdateDoctorProfile(doctor)
	if err != nil {
		return doctorEntities.Doctor{}, err
	}
	return updatedDoctor, nil
}

func (u *DoctorUseCase) GetDetailProfile(doctorID uint) (doctorEntities.Doctor, error) {
	doctor, err := u.doctorRepository.GetDetailProfile(doctorID)
	if err != nil {
		return doctorEntities.Doctor{}, err
	}

	rating, err := u.ratingRepository.GetSummaryRating(doctor.ID)
	if err != nil {
		return doctorEntities.Doctor{}, err
	}

	doctor.RatingPrecentage = (rating.Average / 5) * 100

	return doctor, nil
}