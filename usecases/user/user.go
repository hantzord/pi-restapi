package user

import (
	"capstone/constants"
	userEntitites "capstone/entities/user"
	"capstone/middlewares"
	"capstone/utilities"
	"context"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
	myoauth "golang.org/x/oauth2"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type UserUseCase struct {
	repository    userEntitites.RepositoryInterface
	oauthConfig   *myoauth.Config
	oauthConfigFB *myoauth.Config
}

func NewUserUseCase(repository userEntitites.RepositoryInterface, oauthConfig *myoauth.Config, oauthConfigFB *myoauth.Config) *UserUseCase {
	return &UserUseCase{
		repository:    repository,
		oauthConfig:   oauthConfig,
		oauthConfigFB: oauthConfigFB,
	}
}

func (userUseCase *UserUseCase) Register(user *userEntitites.User) (userEntitites.User, error) {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return userEntitites.User{}, constants.ErrEmptyInputUser
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return userEntitites.User{}, constants.ErrHashedPassword
	}

	user.Password = string(hashedPassword)

	var kode int64
	userResult, kode, err := userUseCase.repository.Register(user)
	if err != nil {
		return userEntitites.User{}, constants.ErrInsertDatabase
	}

	if kode == 1 {
		return userEntitites.User{}, constants.ErrUsernameAlreadyExist
	}

	if kode == 2 {
		return userEntitites.User{}, constants.ErrEmailAlreadyExist
	}

	token, _ := middlewares.CreateToken(userResult.Id)
	userResult.Token = token

	return userResult, nil
}

func (userUseCase *UserUseCase) Login(user *userEntitites.User) (userEntitites.User, error) {
	if user.Username == "" || user.Password == "" {
		return userEntitites.User{}, constants.ErrEmptyInputLogin
	}

	userResult, err := userUseCase.repository.Login(user)
	if err != nil {
		return userEntitites.User{}, constants.ErrUserNotFound
	}

	token, _ := middlewares.CreateToken(userResult.Id)
	userResult.Token = token

	return userResult, nil
}

func (u *UserUseCase) HandleGoogleLogin() string {
	return u.oauthConfig.AuthCodeURL("state-token", myoauth.AccessTypeOffline)
}

func (u *UserUseCase) HandleGoogleCallback(ctx context.Context, code string) (userEntitites.User, error) {
	token, err := u.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return userEntitites.User{}, constants.ErrExcange
	}

	// Membuat layanan OAuth2
	oauth2Service, err := oauth2.NewService(ctx, option.WithTokenSource(u.oauthConfig.TokenSource(ctx, token)))
	if err != nil {
		return userEntitites.User{}, constants.ErrNewServiceGoogle
	}

	userInfoService := oauth2.NewUserinfoV2MeService(oauth2Service)
	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		return userEntitites.User{}, constants.ErrNewUserInfo
	}

	// Cek apakah pengguna sudah ada di database
	result, myCode, err := u.repository.OauthFindByEmail(userInfo.Email)

	if err != nil && myCode == 0 {
		username := utilities.GetFirstNameWithNumbers(userInfo.Name)

		newUser, err := u.repository.Create(userInfo.Email, userInfo.Picture, userInfo.Name, username)
		if err != nil {
			return userEntitites.User{}, constants.ErrInsertOAuth
		}

		tokenJWT, _ := middlewares.CreateToken(newUser.Id)
		newUser.Token = tokenJWT

		return newUser, nil
	}

	if err != nil && myCode == 1 {
		return userEntitites.User{}, err
	}

	tokenJWT, _ := middlewares.CreateToken(result.Id)
	result.Token = tokenJWT

	return result, nil
}

func (u *UserUseCase) GetPointsByUserId(id int) (int, error) {
	result, err := u.repository.GetPointsByUserId(id)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (u *UserUseCase) ResetPassword(email string, password string) error {
	if password == "" {
		return constants.ErrEmptyResetPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = u.repository.ResetPassword(email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) HandleFacebookLogin() string {
	return u.oauthConfigFB.AuthCodeURL("state-token", myoauth.AccessTypeOffline)
}

func (u *UserUseCase) HandleFacebookCallback(ctx context.Context, code string) (userEntitites.User, error) {
	token, err := u.oauthConfigFB.Exchange(ctx, code)
	if err != nil {
		return userEntitites.User{}, constants.ErrExcange
	}

	client := u.oauthConfigFB.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return userEntitites.User{}, constants.ErrNewUserInfo
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
		return userEntitites.User{}, err
	}

	// Cek apakah pengguna sudah ada di database
	result, myCode, err := u.repository.OauthFindByEmail(fbUser.Email)

	if err != nil && myCode == 0 {
		username := utilities.GetFirstNameWithNumbers(fbUser.Name)

		newUser, err := u.repository.Create(fbUser.Email, fbUser.Picture.Data.URL, fbUser.Name, username)
		if err != nil {
			return userEntitites.User{}, constants.ErrInsertOAuth
		}

		tokenJWT, _ := middlewares.CreateToken(newUser.Id)
		newUser.Token = tokenJWT

		return newUser, nil
	}

	if err != nil && myCode == 1 {
		return userEntitites.User{}, err
	}

	tokenJWT, _ := middlewares.CreateToken(result.Id)
	result.Token = tokenJWT

	return result, nil
}

func (u *UserUseCase) UpdateUserProfile(user *userEntitites.User) (userEntitites.User, error) {
	if user.Id == 0 {
		return userEntitites.User{}, constants.ErrUserNotFound
	}

	updatedUser, err := u.repository.UpdateUserProfile(user)
	if err != nil {
		return userEntitites.User{}, err
	}

	return updatedUser, nil
}

func (u *UserUseCase) ChangePassword(userId int, oldPassword, newPassword string) error {
	err := u.repository.ChangePassword(userId, oldPassword, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) UpdateSuccessPointByUserID(id int, pointSpend int) error {
	if pointSpend < 0 {
		return constants.ErrPointSpend
	}

	currentPoints, err := u.repository.GetPointsByUserId(id)
	if err != nil {
		return err
	}

	newPoints := currentPoints - pointSpend
	if newPoints < 0 {
		return constants.ErrInsufficientPoint
	}

	err = u.repository.UpdatePointsByUserID(id, newPoints)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) UpdateFailedPointByUserID(id int, pointSpend int) error {
	if pointSpend < 0 {
		return constants.ErrPointSpend
	}

	currentPoints, err := u.repository.GetPointsByUserId(id)
	if err != nil {
		return err
	}

	newPoints := currentPoints + pointSpend
	if newPoints < 0 {
		return constants.ErrInsufficientPoint
	}

	err = u.repository.UpdatePointsByUserID(id, newPoints)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) ChangeEmail(userId int, newEmail string) error {
	if newEmail == "" {
		return constants.ErrEmptyNewEmail
	}

	err := u.repository.ChangeEmail(userId, newEmail)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) GetDetailedProfile(id int) (userEntitites.User, error) {
	user, err := u.repository.GetDetailedProfile(id)
	if err != nil {
		return userEntitites.User{}, err
	}

	return user, nil
}